package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lumosolutions/yerd/server/internal/config"
	"github.com/lumosolutions/yerd/server/internal/constants"
	"github.com/lumosolutions/yerd/server/internal/core/http/output"
	"github.com/lumosolutions/yerd/server/internal/core/services"
	"github.com/lumosolutions/yerd/server/internal/core/validation"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

type PhpInstaller struct {
	Info          *config.PhpInfo
	Installed     bool
	ReplaceConfig bool
	DownloadUrl   string
	DownloadPath  string
	SrcPath       string
	Output        *output.StreamWriter
	Log           *services.Logger
	Dm            *DependencyManager
	Cm            *CommandManager
	Compile       *CompileManager
}

func NewPhpManager(version string, sw *output.StreamWriter, logger *services.Logger) (*PhpInstaller, error) {
	dm, err := NewDependencyManager(logger)
	if err != nil {
		return nil, err
	}

	cm := NewCommandManager(logger)
	compile := NewCompileManager(logger)

	if info, exists := validation.IsPhpInstalled(version); exists {
		return &PhpInstaller{
			Info:      info,
			Installed: true,
			Output:    sw,
			Log:       logger,
			Dm:        dm,
			Cm:        cm,
			Compile:   compile,
		}, nil
	}

	info, err := GeneratePhpInfo(version)
	if err != nil {
		return nil, err
	}

	return &PhpInstaller{
		Info:      info,
		Installed: false,
		Output:    sw,
		Log:       logger,
		Dm:        dm,
		Cm:        cm,
		Compile:   compile,
	}, nil
}

func (i *PhpInstaller) Install() error {
	if i.Installed {
		return fmt.Errorf("php %s already installed", i.Info.Version)
	}

	i.ReplaceConfig = true

	i.Output.WriteInfo("Installing php %s with extensions", i.Info.Version)
	i.Output.WriteInfo("%s", strings.Join(i.Info.Extensions, ", "))
	i.Output.WriteInfo("Distro: %s", i.Dm.Distro)
	i.Output.WriteInfo("Package Manager: %s", i.Dm.Pm)

	err := utils.Run(
		func() error { return i.installDependencies() },
		func() error { return i.getLatestVersionInfo() },
		func() error { return i.downloadPhp() },
		func() error { return i.compilePhp() },
		func() error { return i.makePhp() },
		func() error { return i.actualInstallPhp() },
		func() error { return i.createDefaultConfigs() },
		func() error { return i.installPeclExtensions() },
		func() error { return i.createSymLinks() },
		func() error { return i.createSystemdService() },
		func() error { return i.storeConfig() },
	)

	if err != nil {
		i.Log.Error("install", err)
		i.Output.WriteError(
			"Failed to install PHP %s: %v",
			i.Info.Version,
			err,
		)
		return err
	}

	return nil
}

func (i *PhpInstaller) storeConfig() error {
	appConfig, _ := config.GetConfig()
	if appConfig.Php == nil {
		appConfig.Php = make(config.PhpConfig)
	}

	appConfig.Php[i.Info.Version] = i.Info
	if err := config.WriteConfig(); err != nil {
		i.Log.Error("store-config", err)
		return err
	}

	return nil
}

func (i *PhpInstaller) createSystemdService() error {
	if !i.ReplaceConfig {
		return nil
	}

	sm, err := CreatePhpSystemdService(i.Info, i.Log)
	if err != nil {
		return err
	}

	if err := sm.Start(); err != nil {
		return err
	}

	return nil
}

func (i *PhpInstaller) createDefaultConfigs() error {
	phpFolder := fmt.Sprintf("php%s", i.Info.Version)
	configDir := filepath.Join(utils.GetYerdEtcPath(), phpFolder)
	iniPath := filepath.Join(configDir, "php.ini")
	fpmPoolConf := filepath.Join(configDir, "php-fpm.d", "www.conf")
	phpFpmConf := filepath.Join(configDir, "php-fpm.conf")
	logFileName := fmt.Sprintf("php%s-fpm.log", i.Info.Version)

	updateIni := i.shouldReplaceConfig(iniPath)
	updatePhpFpmConf := i.shouldReplaceConfig(phpFpmConf)
	updateFpmPoolConf := i.shouldReplaceConfig(fpmPoolConf)

	if updateIni {
		if err := i.updateConfig("php", "php.ini", iniPath, utils.TemplateData{}); err != nil {
			return err
		}
	}

	if updatePhpFpmConf {

		data := utils.TemplateData{
			"pid_path": i.Info.FpmPidLocation,
			"log_path": filepath.Join(utils.GetYerdPhpPath(), constants.PhpLogsPathRelative, logFileName),
			"pool_dir": filepath.Dir(i.Info.PoolConfig),
		}
		if err := i.updateConfig("php", "php-fpm.conf", phpFpmConf, data); err != nil {
			return err
		}
	}

	if updateFpmPoolConf {
		user, _ := utils.GetUser()
		data := utils.TemplateData{
			"version":   i.Info.Version,
			"sock_path": i.Info.FpmSocket,
			"log_path":  filepath.Join(utils.GetYerdPhpPath(), constants.PhpLogsPathRelative, logFileName),
			"user":      user.Username,
			"group":     user.GroupName,
		}
		if err := i.updateConfig("php", "www.conf", fpmPoolConf, data); err != nil {
			return err
		}
	}

	return nil
}

func (i *PhpInstaller) updateConfig(folder, file, path string, data utils.TemplateData) error {
	dm := NewDownloadManager(i.Log, 30*time.Second)

	content, err := dm.FetchFromGitHub(folder, file)
	if err != nil {
		i.Output.WriteError("Failed to download %s", file)
		return err
	}

	fullContent := utils.Template(content, data)
	if err := utils.WriteStringToFile(path, fullContent, constants.FilePermissions); err != nil {
		i.Log.Error("dl", err)
		i.Output.WriteError("Failed to write to %s", file)
		return err
	}

	return nil
}

func (i *PhpInstaller) shouldReplaceConfig(path string) bool {
	if utils.FileExists(path) && !i.ReplaceConfig {
		i.Output.WriteInfo("%s exists, leaving unchanged", filepath.Base(path))
		return false
	}

	return true
}

func (i *PhpInstaller) createSymLinks() error {
	globalBinDir := utils.GetGlobalBinPath()
	localBinDir := utils.GetYerdBinPath()

	name := fmt.Sprintf("php%s", i.Info.Version)
	globalBinPath := filepath.Join(globalBinDir, name)
	localBinPath := filepath.Join(localBinDir, name)
	installPath := filepath.Join(i.Info.InstallPath, "bin", "php")

	if err := utils.EnsureFoldersCreated([]string{globalBinDir, localBinDir}); err != nil {
		return err
	}

	if err := utils.CreateSymlink(installPath, localBinPath); err != nil {
		i.Output.WriteError("Unable to create local Symlink")
		i.Log.Error("symlink", err)
		return err
	}

	if err := utils.CreateSymlink(localBinPath, globalBinPath); err != nil {
		i.Output.WriteError("Unable to create local Symlink")
		i.Log.Error("symlink", err)
		return err
	}

	i.Output.WriteInfo("Symlink created: %s", localBinPath)
	i.Output.WriteInfo("Symlink created: %s", globalBinPath)

	return nil
}

func (i *PhpInstaller) installPeclExtensions() error {
	phpVersion := fmt.Sprintf("php%s", i.Info.Version)
	peclPath := filepath.Join(utils.GetYerdPhpPath(), phpVersion, "bin", "pecl")
	iniDir := filepath.Join(utils.GetYerdPhpPath(), phpVersion, "conf.d")

	for _, extName := range i.Info.Extensions {
		if ext, exists := constants.GetExtension(extName); exists && ext.IsPECL {
			if i.isExtensionLoaded(extName) {
				i.Output.WriteInfo("Extension %s already loaded", extName)
				continue
			}

			_, success := i.Cm.ExecuteCommand(peclPath, "install", ext.PECLName)
			if !success {
				i.Output.WriteError("Failed to install extension %s", extName)
				return fmt.Errorf("failed to install PECL extension %s", extName)
			}

			iniPath := filepath.Join(iniDir, fmt.Sprintf("%s.ini", extName))
			content := fmt.Sprintf("extension=%s.so\n", extName)
			if err := utils.WriteStringToFile(iniPath, content, constants.FilePermissions); err != nil {
				i.Output.WriteError("Failed to created ini file for %s", extName)
				return fmt.Errorf("failed to create ini file for %s: %v", extName, err)
			}
			i.Output.WriteInfo("Created extension %s", extName)
		}
	}

	return nil
}

func (i *PhpInstaller) isExtensionLoaded(extName string) bool {
	phpVersion := fmt.Sprintf("php%s", i.Info.Version)
	phpBin := filepath.Join(utils.GetYerdPhpPath(), phpVersion, "bin", "php")
	response, success := i.Cm.ExecuteCommand(phpBin, "-m")
	if !success {
		return false
	}

	modules := strings.ToLower(response)
	return strings.Contains(modules, strings.ToLower(extName))
}

func (i *PhpInstaller) actualInstallPhp() error {
	phpDir := fmt.Sprintf("php%s", i.Info.Version)
	folders := []string{
		filepath.Join(utils.GetYerdPhpPath(), phpDir),
		filepath.Join(utils.GetWorkingPath(), constants.BinPathRelative),
		filepath.Join(utils.GetYerdPhpPath(), constants.PhpRunPathRelative),
		filepath.Join(utils.GetYerdPhpPath(), constants.PhpLogsPathRelative),
		filepath.Join(utils.GetYerdEtcPath(), phpDir, "conf.d"),
		filepath.Join(utils.GetYerdEtcPath(), phpDir, "php-fpm.d"),
	}

	if err := utils.EnsureFoldersCreated(folders); err != nil {
		return err
	}

	_, success := i.Cm.ExecuteCommandInDir(i.SrcPath, "make", "install")
	if !success {
		i.Log.Info("Install", "Failed to install PHP")
		return fmt.Errorf("unable to install php")
	}

	i.Output.WriteSuccess("PHP %s installed", i.Info.Version)

	utils.RemoveDirectory(i.SrcPath)

	return nil
}

func (i *PhpInstaller) makePhp() error {
	i.Output.WriteInfo("Compiling PHP. This may take a few minutes...")
	if err := i.Compile.Make(i.SrcPath); err != nil {
		return err
	}

	i.Output.WriteSuccess("PHP %s Compiled.", i.Info.ExactVersion)
	return nil
}

func (i *PhpInstaller) compilePhp() error {
	if err := i.Compile.Configure(i.SrcPath, i.getConfigureFlags()); err != nil {
		return err
	}

	i.Output.WriteSuccess("PHP Configured Successfully")

	return nil
}

func (i *PhpInstaller) downloadPhp() error {
	dm := NewDownloadManager(i.Log, 30*time.Second)

	i.Output.WriteInfo("Downloading PHP %s from php.net", i.Info.ExactVersion)

	tempPath := os.TempDir()
	workingPath := utils.GetWorkingPath()
	archivePath := filepath.Join(tempPath, fmt.Sprintf("php-%s.tar.gz", i.Info.ExactVersion))
	extractDir := filepath.Join(tempPath, fmt.Sprintf("php-%s-extract", i.Info.ExactVersion))
	finalExtractDir := filepath.Join(extractDir, fmt.Sprintf("php-%s", i.Info.ExactVersion))
	sourceDir := filepath.Join(workingPath, constants.SrcPathRelative, fmt.Sprintf("php-%s", i.Info.ExactVersion))

	if err := dm.DownloadFile(i.DownloadUrl, archivePath); err != nil {
		return err
	}

	if err := i.Cm.ExtractArchive(archivePath, extractDir); err != nil {
		return err
	}

	utils.ReplaceDirectory(sourceDir)
	if err := i.Cm.CopyRecursive(finalExtractDir, sourceDir); err != nil {
		return err
	}

	i.Output.WriteSuccess("Downloaded PHP %s", i.Info.ExactVersion)
	i.Output.WriteSuccess("Extracted PHP %s", i.Info.ExactVersion)

	i.SrcPath = sourceDir

	utils.RemoveDirectory(extractDir)
	utils.RemoveFile(archivePath)

	return nil
}

func (i *PhpInstaller) getLatestVersionInfo() error {
	i.Output.WriteInfo("Getting Latest Version for PHP %s", i.Info.Version)

	version, download, err := i.FetchLatestMinorVersion()
	if err != nil {
		i.Log.Error("php-version", err)
		return err
	}

	i.Info.ExactVersion = version
	i.DownloadUrl = download

	i.Log.Info("php-version", "Php Version: %s", version)
	i.Log.Info("php-version", "Download URL: %s", download)

	i.Output.WriteSuccess(
		"Latest version for PHP %s is %s",
		i.Info.Version,
		version,
	)

	return nil
}

func (i *PhpInstaller) installDependencies() error {
	deps := constants.GetExtensionDependencies(i.Info.Extensions)
	deps = append(deps, "buildtools", "autoconf", "pkgconfig", "re2c", "oniguruma", "xml", "sqlite")

	if err := i.Dm.InstallPackages(deps); err != nil {
		i.Output.WriteError("Unable to install dependencies: %s", err.Error())
		return err
	}

	return nil
}

type PhpReleaseResponse struct {
	Version           string       `json:"version"`
	Date              string       `json:"date"`
	Tags              []string     `json:"tags"`
	Source            []SourceFile `json:"source"`
	SupportedVersions []string     `json:"supported_versions"`
}

type SourceFile struct {
	Filename string `json:"filename"`
	Name     string `json:"name"`
	SHA256   string `json:"sha256"`
	Date     string `json:"date"`
}

func (i *PhpInstaller) FetchLatestMinorVersion() (string, string, error) {
	phpReleaseUrl := "https://www.php.net/releases/index.php?json&version="
	httpTimeout := 10 * time.Second

	downloadManager := NewDownloadManager(i.Log, httpTimeout)

	var release PhpReleaseResponse
	err := downloadManager.FetchJson(phpReleaseUrl+i.Info.Version, &release)
	if err != nil {
		return "", "", err
	}

	downloadURL := ""
	for _, source := range release.Source {
		if strings.HasSuffix(source.Filename, ".tar.gz") {
			downloadURL = fmt.Sprintf("https://www.php.net/distributions/%s", source.Filename)
			break
		}
	}

	return release.Version, downloadURL, nil
}

func (i *PhpInstaller) getConfigureFlags() []string {
	// TODO: Add Elevated User
	userCtx, _ := utils.GetUser()

	baseFlags := []string{
		fmt.Sprintf("--prefix=%s/php%s", utils.GetYerdPhpPath(), i.Info.Version),
		fmt.Sprintf("--with-config-file-path=%s/php%s", utils.GetYerdEtcPath(), i.Info.Version),
		fmt.Sprintf("--with-config-file-scan-dir=%s/php%s/conf.d", utils.GetYerdEtcPath(), i.Info.Version),
		"--enable-fpm",
		fmt.Sprintf("--with-fpm-user=%s", userCtx.Username),
		fmt.Sprintf("--with-fpm-group=%s", userCtx.GroupName),
		"--enable-cli",
		"--with-pear",
	}

	extensionFlags := constants.GetExtensionConfigureFlags(i.Info.Extensions)
	return append(baseFlags, extensionFlags...)
}

func GeneratePhpInfo(version string) (*config.PhpInfo, error) {
	if !validation.IsPhpVersionValid(version) {
		return nil, fmt.Errorf("php version %s is not valid", version)
	}

	workingPath := utils.GetWorkingPath()
	installPath := filepath.Join(workingPath, constants.PhpPathRelative, "php"+version)
	runPath := filepath.Join(workingPath, constants.PhpPathRelative, "run")
	pidPath := filepath.Join(runPath, "php"+version+"-fpm.pid")
	socketPath := filepath.Join(runPath, "php"+version+"-fpm.sock")
	phpEtcPath := filepath.Join(workingPath, constants.EtcPathRelative, "php"+version)
	iniPath := filepath.Join(phpEtcPath, "php.ini")
	fpmConfPath := filepath.Join(phpEtcPath, "php-fpm.conf")
	poolConfPath := filepath.Join(phpEtcPath, "php-fpm.d", "php-pool.conf")
	peclPath := filepath.Join(phpEtcPath, "conf.d")

	return &config.PhpInfo{
		Version:        version,
		ExactVersion:   "",
		InstallPath:    installPath,
		IsCLI:          false,
		Global:         utils.IsRunningElevated(),
		Extensions:     constants.DefaultExtensions,
		NeedsRebuild:   false,
		FpmPidLocation: pidPath,
		FpmSocket:      socketPath,
		PhpIniLocation: iniPath,
		FpmConfig:      fpmConfPath,
		PoolConfig:     poolConfPath,
		PeclPath:       peclPath,
	}, nil
}
