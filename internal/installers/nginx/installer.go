package nginx

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
)

type NginxInstaller struct {
	Info       *constants.NginxConfig
	IsUpdate   bool
	Spinner    *utils.Spinner
	DepManager *manager.DependencyManager
}

func NewNginxInstaller(update bool) (*NginxInstaller, error) {
	s := utils.NewSpinner("Starting Nginx Installer...")
	s.SetDelay(150)

	depMan, err := manager.NewDependencyManager()
	if err != nil {
		s.AddErrorStatus("Failed to create a new dependency manager")
		s.StopWithError("No action taken")
		return nil, err
	}

	return &NginxInstaller{
		Info:       constants.GetNginxConfig(),
		IsUpdate:   update,
		Spinner:    s,
		DepManager: depMan,
	}, nil
}

func (installer *NginxInstaller) Install() error {
	installer.Spinner.Start()

	err := utils.RunAll(
		func() error { return installer.installDependencies() },
		func() error { return installer.prepareInstall() },
		func() error { return installer.downloadSource() },
		func() error { return installer.compileAndInstall() },
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
