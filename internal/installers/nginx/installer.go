package nginx

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
)

type NginxInstaller struct {
	Info        *constants.NginxConfig
	IsUpdate    bool
	ForceConfig bool
	Spinner     *utils.Spinner
	DepManager  *manager.DependencyManager
}

func NewNginxInstaller(update, forceConfig bool) (*NginxInstaller, error) {
	s := utils.NewSpinner("Starting Nginx Installer...")
	s.SetDelay(150)

	depMan, err := manager.NewDependencyManager()
	if err != nil {
		s.AddErrorStatus("Failed to create a new dependency manager")
		s.StopWithError("No action taken")
		return nil, err
	}

	return &NginxInstaller{
		Info:        constants.GetNginxConfig(),
		IsUpdate:    update,
		ForceConfig: forceConfig,
		Spinner:     s,
		DepManager:  depMan,
	}, nil
}

func (installer *NginxInstaller) Install() error {
	installer.Spinner.Start()

	err := utils.RunAll(
		// func() error { return installer.installDependencies() },
		// func() error { return installer.prepareInstall() },
		// func() error { return installer.downloadSource() },
		// func() error { return installer.compileAndInstall() },
		// func() error { return installer.addNginxConf() },
		// func() error { return installer.addSystemdService() },
		// func() error { return installer.writeConfig() },
		func() error { return installer.createCerts() },
	)

	if err != nil {
		return err
	}

	installer.Spinner.StopWithSuccess("Installed Successfully")

	return nil
}

func (installer *NginxInstaller) installDependencies() error {
	installer.Spinner.UpdatePhrase("Installing Dependencies...")
	if err := installer.DepManager.InstallWebDependencies(); err != nil {
		installer.Spinner.StopWithError("Failed to install dependencies")
		return err
	}

	installer.Spinner.AddSuccessStatus("Dependencies Installed")
	return nil
}

func (installer *NginxInstaller) prepareInstall() error {
	installer.Spinner.UpdatePhrase("Created required folders")
	requiredDirs := []string{
		installer.Info.ConfigPath,
		installer.Info.LogPath,
		installer.Info.RunPath,
		installer.Info.TempPath,
		installer.Info.SourcePath,
		filepath.Join(installer.Info.InstallPath, "sbin"),
		filepath.Join(constants.YerdWebDir, "nginx", "sites-available"),
	}

	for _, dir := range requiredDirs {
		if err := utils.CreateDirectory(dir); err != nil {
			installer.Spinner.AddErrorStatus("Failed to create directory: %s", dir)
			installer.Spinner.StopWithError("Installation stopped failure in setup")
			return err
		}
	}

	installer.Spinner.AddSuccessStatus("Directories created successfully")

	return nil
}

func (installer *NginxInstaller) downloadSource() error {
	installer.Spinner.UpdatePhrase("Downloading Nginx")
	archivePath := filepath.Join(os.TempDir(), "nginx.tar.gz")
	if err := utils.DownloadFile(installer.Info.DownloadURL, archivePath, nil); err != nil {
		installer.Spinner.AddErrorStatus("Unable to download Nginx")
		installer.Spinner.AddInfoStatus("- Error: %v", err)
		installer.Spinner.StopWithError("Failed to download Nginx")
		return err
	}

	userCtx, err := utils.GetRealUser()
	if err != nil {
		installer.Spinner.StopWithError("Failed to identify real user")
		return err
	}

	if err := utils.ExtractArchive(archivePath, installer.Info.SourcePath, userCtx); err != nil {
		installer.Spinner.StopWithError("Failed to extract Nginx")
		return err
	}

	installer.Spinner.AddSuccessStatus("Downloaded nginx successfully")

	return nil
}

func (installer *NginxInstaller) compileAndInstall() error {
	installer.Spinner.UpdatePhrase("Compiling Nginx...")

	buildPath := filepath.Join(installer.Info.SourcePath, fmt.Sprintf("nginx-%s", installer.Info.Version))

	if !utils.FileExists(filepath.Join(buildPath, "/configure")) {
		installer.Spinner.StopWithError("No configure script for Nginx")
		return fmt.Errorf("configure script not found in source directory")
	}

	_, success := utils.ExecuteCommandInDir(
		buildPath,
		"./configure",
		installer.Info.BuildFlags...,
	)

	if !success {
		installer.Spinner.StopWithError("Unable to configure Nginx")
		return fmt.Errorf("unable to configure nginx")
	}

	installer.Spinner.AddSuccessStatus("Nginx Configured Successfully")
	installer.Spinner.UpdatePhrase("Installing Nginx...")

	_, success = utils.ExecuteCommandInDir(
		buildPath,
		"make",
		"install",
	)

	if !success {
		installer.Spinner.StopWithError("Unable to install Nginx")
		return fmt.Errorf("unable to install nginx")
	}

	utils.RemoveFolder(installer.Info.SourcePath)

	return nil
}

func (installer *NginxInstaller) addNginxConf() error {
	if installer.IsUpdate && !installer.ForceConfig {
		installer.Spinner.AddInfoStatus("- Not updating nginx.conf")
		return nil
	}

	installer.Spinner.UpdatePhrase("Downloading nginx.conf")

	content, err := utils.FetchFromGitHub("nginx", "nginx.conf")
	if err != nil {
		utils.LogError(err, "addConf")
		installer.Spinner.AddErrorStatus("Failed to download nginx configuration")
		installer.Spinner.StopWithError("nginx.conf download failed")
		return err
	}

	content = utils.Template(content, utils.TemplateData{
		"user": "root",
	})

	filePath := filepath.Join(installer.Info.ConfigPath, "nginx.conf")

	err = utils.WriteStringToFile(filePath, content, constants.FilePermissions)
	if err != nil {
		utils.LogError(err, "addConf")
		installer.Spinner.StopWithError("Failed to write to nginx.conf")
		return err
	}

	installer.Spinner.AddInfoStatus("- Downloaded Nginx.conf")
	installer.Spinner.AddInfoStatus("- Stored Nginx.conf")
	installer.Spinner.AddSuccessStatus("Nginx Configured Successfully")
	return nil
}

func (installer *NginxInstaller) addSystemdService() error {
	if installer.IsUpdate && !installer.ForceConfig {
		installer.Spinner.AddInfoStatus("- Not updating systemd config")
		return nil
	}

	installer.Spinner.UpdatePhrase("Configuring Systemd")

	systemdPath := filepath.Join(constants.SystemdDir, "yerd-nginx.service")
	content, err := utils.FetchFromGitHub("nginx", "systemd.conf")
	if err != nil {
		utils.LogError(err, "systemd")
		installer.Spinner.AddErrorStatus("Failed to download systemd configuration")
		installer.Spinner.StopWithError("systemd.conf download failed")
		return err
	}

	utils.WriteStringToFile(systemdPath, content, constants.FilePermissions)

	installer.Spinner.AddInfoStatus("Created %s", filepath.Base(systemdPath))

	if err := utils.SystemdReload(); err != nil {
		utils.LogInfo("setupSystemd", "Unable to reload daemons")
		return err
	}

	installer.Spinner.AddInfoStatus("[Systemd] Reloaded daemons")

	serviceName := "yerd-nginx"
	utils.SystemdStopService(serviceName)
	if err := utils.SystemdStartService(serviceName); err != nil {
		utils.LogInfo("setupSystemd", "Unable to start service %s", serviceName)
		installer.Spinner.StopWithError("Unable to start service %s", serviceName)
		return fmt.Errorf("unable to start service %s", serviceName)
	}

	installer.Spinner.AddInfoStatus("[Systemd] Started '%s' successfully", serviceName)
	installer.Spinner.AddSuccessStatus("Systemd Configured")

	return nil
}

func (installer *NginxInstaller) writeConfig() error {
	installer.Spinner.UpdatePhrase("Writing YERD Configuration")

	var existing *config.WebConfig
	err := config.GetStruct("web", &existing)
	if err != nil || existing == nil {
		newConfig := config.WebConfig{
			Installed: true,
		}

		config.SetStruct("web", newConfig)
		installer.Spinner.AddSuccessStatus("YERD Configuration Created")
		return nil
	}

	if existing.Installed {
		installer.Spinner.AddInfoStatus("- YERD configuration does not need updating")
		return nil
	}

	existing.Installed = true
	config.SetStruct("web", existing)

	hostManager := utils.NewHostsManager()
	hostManager.Install()

	return nil
}

func (installer *NginxInstaller) createCerts() error {
	installer.Spinner.UpdatePhrase("Generated Root CA Certificate")
	certManager := manager.NewCertificateManager()
	if err := certManager.GenerateCaCertificate("yerd"); err != nil {
		installer.Spinner.AddErrorStatus("Failed to generate CA Certificate")
		return err
	}

	installer.Spinner.AddSuccessStatus("Root CA Certificate Created")

	return nil
}
