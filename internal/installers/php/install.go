package php

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
)

type PhpInstaller struct {
	version         string
	info            *PhpVersionInfo
	useCache        bool
	update          bool
	updateConfig    bool
	useExactVersion bool
	exactVersion    string
	extensions      []string
	spinner         *utils.Spinner
	depManager      *manager.DependencyManager
	installPath     string
	err             error
}

func NewPhpInstaller(version string, useCache, updateConfig bool) (*PhpInstaller, error) {
	if err := utils.EnsureYerdDirectories(); err != nil {
		utils.LogError(err, "php")
		return nil, fmt.Errorf("failed to create yerd directories")
	}

	var install *config.PhpInfo
	var update bool = false
	var extensions = constants.GetDefaultExtensions()
	if err := config.GetStruct(fmt.Sprintf("php.[%s]", version), &install); err == nil {
		update = true
		extensions = install.Extensions
		extensions = utils.AddUnique(extensions, install.AddExtensions...)
		extensions = utils.RemoveItems(extensions, install.RemoveExtensions...)
	}

	s := utils.NewSpinner("Starting Installer...")
	s.SetDelay(150)

	return &PhpInstaller{
		version:      version,
		spinner:      s,
		update:       update,
		useCache:     useCache,
		updateConfig: updateConfig,
		extensions:   extensions,
	}, nil
}

func (installer *PhpInstaller) Install() error {
	displayType := "Installing"
	if installer.update {
		displayType = "Rebuilding"
	}

	fmt.Println(displayType + " PHP " + installer.version + " with extensions")
	utils.PrintExtensionsGrid(installer.extensions)
	fmt.Println()

	installer.spinner.Start()

	installer.
		run(installer.identifySystem).
		run(installer.getVersionInfo).
		run(installer.conflictingBinaries).
		run(installer.installDeps).
		run(installer.downloadPhp).
		run(installer.configurePhp).
		run(installer.compilePhp).
		run(installer.makePhp).
		run(installer.createSymlinks).
		run(installer.verifyInstall).
		run(installer.createDefaultConfig).
		run(installer.setupSystemdService).
		run(installer.writeConfig)

	if installer.err != nil {
		return installer.err
	}

	installer.spinner.StopWithSuccess("PHP %s Installed", installer.version)

	return nil
}

func (installer *PhpInstaller) run(fn func() error) *PhpInstaller {
	if installer.err == nil {
		installer.err = fn()
	}

	return installer
}

func (installer *PhpInstaller) identifySystem() error {
	installer.spinner.UpdatePhrase("Identifying System...")

	manager, err := manager.NewDependencyManager()
	if err != nil {
		installer.spinner.StopWithError("Unable to identify the system")
		return fmt.Errorf("failed to create dependency manager")
	}

	installer.depManager = manager

	installer.spinner.AddSuccessStatus("System Identification Complete")
	installer.spinner.AddInfoStatus("Distro: %s", installer.depManager.GetDistro())
	installer.spinner.AddInfoStatus("Package Manager: %s", installer.depManager.GetPackageManager())

	return nil
}

func (installer *PhpInstaller) installDeps() error {
	installer.spinner.UpdatePhrase("Installing Dependencies...")

	if err := installer.depManager.InstallBuildDependencies(); err != nil {
		installer.spinner.StopWithError("Failed to install build dependencies")
		return err
	}

	installer.spinner.AddSuccessStatus("Installed Build Dependencies")

	if err := installer.depManager.InstallExtensionDependencies(installer.extensions); err != nil {
		installer.spinner.StopWithError("Failed to install extension dependencies")
		return err
	}

	installer.spinner.AddSuccessStatus("Installed Extension Dependencies")
	return nil
}

func (installer *PhpInstaller) downloadPhp() error {
	installer.spinner.UpdatePhrase("Downloading Source from php.net...")

	userCtx, err := utils.GetRealUser()
	if err != nil {
		utils.LogError(err, "extract")
		installer.spinner.StopWithError("Unable to get user information")
		return fmt.Errorf("error getting user information")
	}

	if err := utils.DownloadFile(installer.info.DownloadURL, installer.info.ArchivePath, nil); err != nil {
		utils.LogError(err, "download")
		installer.spinner.StopWithError("Unable to download")
		PrintVersionFetchError(installer.version)
		return fmt.Errorf("unable to download php%s", installer.info.MajorMinor)
	}

	if err := utils.ExtractArchive(installer.info.ArchivePath, installer.info.ExtractPath, userCtx); err != nil {
		utils.LogError(err, "download")
		os.Remove(installer.info.ArchivePath)
		installer.spinner.StopWithError("Unable to extract")
		PrintVersionFetchError(installer.version)
		return fmt.Errorf("unable to extract php%s", installer.info.MajorMinor)
	}

	tempSourcePath := filepath.Join(installer.info.ExtractPath, fmt.Sprintf("php-%s", installer.info.Version))
	utils.ReplaceDirectory(installer.info.SourcePath)
	utils.CopyRecursive(tempSourcePath, installer.info.SourcePath)
	utils.ChownRecursive(installer.info.SourcePath, userCtx.UID, userCtx.GID)

	os.Remove(installer.info.ArchivePath)
	os.RemoveAll(installer.info.ExtractPath)

	installer.spinner.AddSuccessStatus("PHP Source Download Complete")

	return nil
}

func (installer *PhpInstaller) getVersionInfo() error {
	if !installer.useExactVersion {
		installer.spinner.UpdatePhrase("Fetching Latest Version...")

		info, err := getPhpVersionInfo(installer.version, installer.useCache)
		if err != nil {
			installer.spinner.StopWithError("Unable to fetch latest version from php.net")
			return err
		}

		installer.info = info

		installer.spinner.AddSuccessStatus("Fetched Latest Version")
		installer.spinner.AddInfoStatus("Version: %s", installer.info.Version)

		return nil
	}

	installer.spinner.UpdatePhrase(fmt.Sprintf("Fetching Version %s...", installer.exactVersion))
	version, download, err := FetchSpecificVersion(installer.exactVersion)
	if err != nil {
		installer.spinner.StopWithError(
			fmt.Sprintf("Unable to fetch PHP version %s from php.net", installer.exactVersion),
		)
		return err
	}

	tempDir := os.TempDir()

	installer.info = &PhpVersionInfo{
		MajorMinor:     installer.version,
		Version:        version,
		DownloadURL:    download,
		ConfigureFlags: getConfigureFlags(installer.version, installer.extensions),
		SourcePackage:  fmt.Sprintf("php-%s", installer.version),
		ArchivePath:    filepath.Join(tempDir, fmt.Sprintf("php-%s.tar.gz", version)),
		ExtractPath:    filepath.Join(tempDir, fmt.Sprintf("php-%s-extract", version)),
		SourcePath:     filepath.Join(constants.YerdPHPDir, "src", fmt.Sprintf("php-%s", installer.version)),
	}

	installer.spinner.AddSuccessStatus("Fetched Specific Version")
	installer.spinner.AddInfoStatus("Version: %s", installer.info.Version)

	return nil
}

func (installer *PhpInstaller) conflictingBinaries() error {
	if isYerdManaged("php") {
		installer.spinner.AddWarningStatus("Conflicting root 'php' installation")
	} else {
		installer.spinner.AddInfoStatus("No conflicting 'php' installations")
	}

	if isYerdManaged(fmt.Sprintf("php%s", installer.version)) {
		installer.spinner.AddErrorStatus("Conflicting 'php%s' installation", installer.version)
		installer.spinner.StopWithError("Unable to installed php%s due to conflicting executable", installer.version)
		return fmt.Errorf("conflicting php executable")
	}

	return nil
}

func isYerdManaged(command string) bool {
	if path, exists := utils.CommandExists(command); exists {
		if utils.IsSymlink(path) {
			rootPath, err := utils.ReadSymlink(path)
			utils.LogInfo("php", "SymLink for '%s' at '%s' is '%s'", command, path, rootPath)
			return err != nil || !strings.HasPrefix(rootPath, constants.YerdBaseDir)
		} else {
			utils.LogInfo("php", "'%s' at '%s' is not a symlink", command, path)
			return !strings.HasPrefix(path, constants.YerdBaseDir)
		}
	}

	return false
}

func (installer *PhpInstaller) configurePhp() error {
	utils.LogInfo("php", "Starting to configure PHP")
	installer.spinner.UpdatePhrase("Configuring PHP...")

	configurePath := filepath.Join(installer.info.SourcePath, "configure")
	if exists := utils.FileExists(configurePath); !exists {
		installer.spinner.StopWithError("Configure script not found")
		return fmt.Errorf("configure script not found")
	}

	if err := utils.Chmod(configurePath, 0755); err != nil {
		installer.spinner.StopWithError("Unable to make configure script executable")
		return fmt.Errorf("unable to make configure executable")
	}

	args := append([]string{"/bin/bash", configurePath}, installer.info.ConfigureFlags...)
	if _, success := utils.ExecuteCommandInDirAsUser(installer.info.SourcePath, args[0], args[1:]...); !success {
		installer.spinner.StopWithError("Unable to run configure script")
		return fmt.Errorf("unable to run configure script")
	}

	utils.LogInfo("php", "Configure Complete")
	installer.spinner.AddSuccessStatus("PHP Configured Successfully")

	return nil
}

func (installer *PhpInstaller) compilePhp() error {
	utils.LogInfo("compile", "Starting to compile PHP")
	installer.spinner.UpdatePhrase("Compiling PHP (may take a few minutes)...")

	nproc := utils.GetProcessorCount()
	if _, success := utils.ExecuteCommandInDirAsUser(installer.info.SourcePath, "make", fmt.Sprintf("-j%d", nproc)); !success {
		installer.spinner.StopWithError("Unable to compile PHP")
		return fmt.Errorf("unable to compile php")
	}

	utils.LogInfo("compile", "Compile complete")
	installer.spinner.AddSuccessStatus("PHP Compiled Successfully")

	return nil
}

func (installer *PhpInstaller) makePhp() error {
	utils.LogInfo("install", "Starting to actually install PHP")
	installer.spinner.UpdatePhrase("Installing PHP...")

	output, success := utils.ExecuteCommandInDir(installer.info.SourcePath, "make", "install")
	if !success {
		installer.spinner.StopWithError("Failed to install PHP")
		utils.LogDebug("install", "%s", output)
		return fmt.Errorf("unable to install php")
	}

	installer.installPath = fmt.Sprintf("%s/php%s", constants.YerdPHPDir, installer.version)

	utils.RemoveFolder(installer.info.SourcePath)

	installer.spinner.AddSuccessStatus("PHP Installed Successfully")

	return nil
}

func (installer *PhpInstaller) createSymlinks() error {
	installer.spinner.UpdatePhrase("Created Symlinks")

	installedBinary := installer.installPath + "/bin/php"
	localBinaryPath := constants.YerdBinDir + "/php" + installer.version
	globalBinaryPath := constants.SystemBinDir + "/php" + installer.version

	if err := utils.CreateSymlink(installedBinary, localBinaryPath); err != nil {
		installer.spinner.StopWithError("Unable to create local Symlink")
		utils.LogError(err, "symlink")
		return err
	}

	if err := utils.CreateSymlink(localBinaryPath, globalBinaryPath); err != nil {
		installer.spinner.StopWithError("Unable to create global Symlink")
		utils.LogError(err, "symlink")
		return err
	}

	installer.spinner.AddSuccessStatus("Created Symlinks")
	installer.spinner.AddInfoStatus("Local Binary: %s", localBinaryPath)
	installer.spinner.AddInfoStatus("Global Binary: %s", globalBinaryPath)

	return nil
}

func (installer *PhpInstaller) verifyInstall() error {
	installer.spinner.UpdatePhrase("Verifying Installation")

	binaries := []string{
		installer.installPath + "/bin/php",
		constants.YerdBinDir + "/php" + installer.version,
		constants.SystemBinDir + "/php" + installer.version,
		"php" + installer.version,
	}

	for _, binary := range binaries {
		utils.LogInfo("verify", "Testing: %s", binary)
		if _, success := utils.ExecuteCommandAsUser(binary, "-v"); !success {
			installer.spinner.StopWithError("PHP binary not executable: %s", binary)
			return fmt.Errorf("binary not executable")
		}

		installer.spinner.AddInfoStatus("Verified %s", binary)
	}

	installer.spinner.AddSuccessStatus("Installation Verified")

	return nil
}

func (installer *PhpInstaller) createDefaultConfig() error {
	installer.spinner.UpdatePhrase("Check/Update Configuraiton")

	configDir := filepath.Join(constants.YerdEtcDir, "php"+installer.version)
	iniPath := filepath.Join(configDir, "php.ini")
	fpmPoolConf := filepath.Join(configDir, constants.FPMPoolDir, constants.FPMPoolConfig)
	phpFpmConf := filepath.Join(configDir, "php-fpm.conf")

	utils.CreateDirectory(filepath.Join(configDir, constants.FPMPoolDir))
	utils.CreateDirectory(filepath.Join(constants.YerdPHPDir, "logs"))
	utils.CreateDirectory(constants.FPMSockDir)
	utils.CreateDirectory(constants.FPMPidDir)

	updateIni := installer.shouldReplaceConfig(iniPath)
	updatePhpFpmConf := installer.shouldReplaceConfig(phpFpmConf)
	updateFpmPoolConf := installer.shouldReplaceConfig(fpmPoolConf)

	if updateIni {
		if err := installer.downloadAndReplace("php", "php.ini", iniPath, utils.TemplateData{}); err != nil {
			return err
		}
	}

	if updatePhpFpmConf {
		data := utils.TemplateData{
			"pid_path": filepath.Join(constants.FPMSockDir, fmt.Sprintf("php%s-fpm.pid", installer.version)),
			"log_path": filepath.Join(constants.FPMLogDir, fmt.Sprintf("php%s-fpm.log", installer.version)),
			"pool_dir": filepath.Join(constants.YerdEtcDir, "php"+installer.version, constants.FPMPoolDir),
		}
		if err := installer.downloadAndReplace("php", "php-fpm.conf", phpFpmConf, data); err != nil {
			return err
		}
	}

	if updateFpmPoolConf {
		data := utils.TemplateData{
			"version":   installer.version,
			"sock_path": filepath.Join(constants.FPMSockDir, fmt.Sprintf("php%s-fpm.sock", installer.version)),
			"log_path":  filepath.Join(constants.FPMLogDir, fmt.Sprintf("php%s-fpm.log", installer.version)),
			"user":      GetFPMUser(),
			"group":     GetFPMGroup(),
		}
		if err := installer.downloadAndReplace("php", "www.conf", fpmPoolConf, data); err != nil {
			return err
		}
	}

	installer.spinner.AddSuccessStatus("PHP Configured Successfully")

	return nil
}

func (installer *PhpInstaller) setupSystemdService() error {
	installer.spinner.UpdatePhrase("Configuring Systemd")

	systemdPath := filepath.Join(constants.SystemdDir, fmt.Sprintf("yerd-php%s-fpm.service", installer.version))
	updateSystemdConf := installer.shouldReplaceConfig(systemdPath)

	if updateSystemdConf {
		phpVersionStr := fmt.Sprintf("php%s", installer.version)
		configDir := filepath.Join(constants.YerdEtcDir, "php"+installer.version)
		phpFpmConf := filepath.Join(configDir, "php-fpm.conf")
		data := utils.TemplateData{
			"version":          installer.version,
			"pid_path":         filepath.Join(constants.FPMSockDir, fmt.Sprintf("%s-fpm.pid", phpVersionStr)),
			"fpm_binary_path":  filepath.Join(constants.YerdPHPDir, phpVersionStr, "sbin", "php-fpm"),
			"main_config_path": phpFpmConf,
		}

		if err := installer.downloadAndReplace("php", "systemd.conf", systemdPath, data); err != nil {
			return err
		}

		installer.spinner.AddInfoStatus("Created %s", filepath.Base(systemdPath))

		if err := utils.SystemdReload(); err != nil {
			utils.LogInfo("setupSystemd", "Unable to reload daemons")
			return err
		}

		installer.spinner.AddInfoStatus("[Systemd] Reloaded daemons")
	}

	serviceName := fmt.Sprintf("yerd-php%s-fpm", installer.version)
	utils.SystemdStopService(serviceName)
	if err := utils.SystemdStartService(serviceName); err != nil {
		utils.LogInfo("setupSystemd", "Unable to start service %s", serviceName)
		return fmt.Errorf("unable to start service %s", serviceName)
	}

	utils.SystemdEnable(serviceName)

	installer.spinner.AddInfoStatus("[Systemd] Started '%s' successfully", serviceName)
	installer.spinner.AddSuccessStatus("Systemd Configured")

	return nil
}

func (installer *PhpInstaller) shouldReplaceConfig(path string) bool {
	if utils.FileExists(path) && !installer.updateConfig {
		utils.LogInfo("phpini", "%s exists, not updating", filepath.Base(path))
		installer.spinner.AddInfoStatus("%s exists, leaving unchanged", filepath.Base(path))
		return false
	}

	return true
}

func (installer *PhpInstaller) downloadAndReplace(folder, file, path string, data utils.TemplateData) error {
	content, err := utils.FetchFromGitHub(folder, file)
	if err != nil {
		installer.spinner.StopWithError("Failed to download %s", file)
		return err
	}

	fullContent := utils.Template(content, data)
	if err := utils.WriteStringToFile(path, fullContent, constants.FilePermissions); err != nil {
		utils.LogError(err, "dl")
		installer.spinner.StopWithError("Failed to write to %s", file)
		return err
	}

	return nil
}

func (installer *PhpInstaller) writeConfig() error {
	var existing *config.PhpInfo
	configPath := fmt.Sprintf("php.[%s]", installer.version)

	if installer.update {
		existing, _ = config.GetInstalledPhpInfo(installer.version)
	}

	isCli := false
	if existing != nil {
		isCli = existing.IsCLI
	}

	data := config.PhpInfo{
		Version:          installer.version,
		InstallPath:      installer.installPath,
		InstallDate:      time.Now(),
		InstalledVersion: installer.info.Version,
		Extensions:       installer.extensions,
		IsCLI:            isCli,
		RemoveExtensions: []string{},
		AddExtensions:    []string{},
	}

	config.SetStruct(configPath, data)

	return nil
}

func (installer *PhpInstaller) UseVersion(fullVersion string) {
	installer.useExactVersion = true
	installer.exactVersion = fullVersion
}

// getConfigureFlags gets the default configure flags and appends the extension
// configure flags required to build a specific version of PHP
// majorMinor: PHP version, for example 8.1
// extensions: The list of PHP extenstions that are to be installed
func getConfigureFlags(majorMinor string, extensions []string) []string {
	baseFlags := []string{
		fmt.Sprintf("--prefix=%s/php%s", constants.YerdPHPDir, majorMinor),
		fmt.Sprintf("--with-config-file-path=%s/php%s", constants.YerdEtcDir, majorMinor),
		fmt.Sprintf("--with-config-file-scan-dir=%s/php%s/conf.d", constants.YerdEtcDir, majorMinor),
		"--enable-fpm",
		fmt.Sprintf("--with-fpm-user=%s", GetFPMUser()),
		fmt.Sprintf("--with-fpm-group=%s", GetFPMGroup()),
		"--enable-cli",
	}

	extensionFlags := constants.GetExtensionConfigureFlags(extensions)

	return append(baseFlags, extensionFlags...)
}

// printVersionFetchError displays helpful error message when version fetching fails.
// version: PHP version that failed to fetch.
func PrintVersionFetchError(version string) {
	fmt.Printf("❌ Could not fetch latest PHP versions from PHP.net\n")
	fmt.Printf("💡 This is required to:\n")
	fmt.Printf("   • Get the latest stable version of PHP %s\n", version)
	fmt.Printf("   • Download the source code from php.net\n")
	fmt.Printf("   • Ensure a secure and up-to-date installation\n\n")
	fmt.Printf("🔍 Troubleshooting:\n")
	fmt.Printf("   • Check your internet connection\n")
	fmt.Printf("   • Verify PHP.net is accessible: curl -I https://www.php.net\n")
	fmt.Printf("   • Try again in a few moments\n")
}

func RunRebuild(data *config.PhpInfo, nocache, config bool) error {
	if nocache {
		fmt.Println("ℹ️  Bypassing cache to get latest version information")
	}

	if config {
		fmt.Println("ℹ️  Recreating configuration if it already exists")
	}

	fmt.Println()

	installer, err := NewPhpInstaller(data.Version, nocache, config)
	if err != nil {
		return err
	}

	installer.UseVersion(data.InstalledVersion)
	if err := installer.Install(); err != nil {
		return err
	}

	fmt.Println()

	fmt.Println("Rebuild has completed...")
	fmt.Println("Thanks for using YERD")

	return nil
}
