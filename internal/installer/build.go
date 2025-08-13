package installer

import (
	"fmt"
	"os"
	"strings"

	"github.com/LumoSolutions/yerd/internal/dependencies"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/pkg/php"
)

func buildAndInstall(versionInfo php.VersionInfo, sourceDir string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Starting build process in: %s", sourceDir)
	
	// Use default extensions for new installations
	defaultExtensions := []string{
		"mbstring", "bcmath", "opcache", "curl", "openssl", "zip", 
		"sockets", "mysqli", "pdo-mysql", "gd", "jpeg", "freetype",
	}
	
	// Install extension dependencies
	depMgr, err := dependencies.NewDependencyManager()
	if err != nil {
		utils.SafeLog(logger, "Failed to initialize dependency manager: %v", err)
		return fmt.Errorf("failed to initialize dependency manager: %v", err)
	}
	
	utils.SafeLog(logger, "Installing extension dependencies for: %v", defaultExtensions)
	if err := depMgr.InstallExtensionDependencies(defaultExtensions); err != nil {
		utils.SafeLog(logger, "Failed to install extension dependencies: %v", err)
		return fmt.Errorf("failed to install extension dependencies: %v", err)
	}
	
	configureFlags := php.GetConfigureFlagsForVersion(versionInfo.Version, defaultExtensions)
	utils.SafeLog(logger, "Configure flags: %v", configureFlags)
	utils.SafeLog(logger, "Extensions: %v", defaultExtensions)
	
	oldDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)
	
	if err := os.Chdir(sourceDir); err != nil {
		utils.SafeLog(logger, "Failed to change to source directory: %v", err)
		return fmt.Errorf("failed to change to source directory: %v", err)
	}
	
	if err := configureSource(configureFlags, logger); err != nil {
		return err
	}
	
	if err := makeSource(logger); err != nil {
		return err
	}
	
	if err := makeInstall(logger); err != nil {
		return err
	}
	
	return nil
}


func configureSource(configureFlags []string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Running ./configure with flags: %v", configureFlags)
	
	spinner := utils.NewLoadingSpinner("Configuring build...")
	spinner.Start()
	
	args := append([]string{"./configure"}, configureFlags...)
	output, err := utils.ExecuteCommand("bash", args...)
	if err != nil {
		spinner.Stop("❌ Configure failed")
		utils.SafeLog(logger, "Configure failed: %v", err)
		utils.SafeLog(logger, "Configure output: %s", output)
		return fmt.Errorf("configure failed (see log for details)")
	}
	
	spinner.Stop("✓ Configure complete")
	utils.SafeLog(logger, "Configure completed successfully")
	
	return nil
}

func makeSource(logger *utils.Logger) error {
	utils.SafeLog(logger, "Running make...")
	
	spinner := utils.NewLoadingSpinner("Building PHP (this may take several minutes)...")
	spinner.Start()
	
	nproc := "4"
	if output, err := utils.ExecuteCommand("nproc"); err == nil && strings.TrimSpace(output) != "" {
		nproc = strings.TrimSpace(output)
	}
	
	output, err := utils.ExecuteCommand("make", "-j"+nproc)
	if err != nil {
		spinner.Stop("❌ Build failed")
		utils.SafeLog(logger, "Build failed: %v", err)
		utils.SafeLog(logger, "Make output: %s", output)
		return fmt.Errorf("build failed (see log for details)")
	}
	
	spinner.Stop("✓ Build complete")
	utils.SafeLog(logger, "Build completed successfully")
	
	return nil
}

func makeInstall(logger *utils.Logger) error {
	utils.SafeLog(logger, "Running make install...")
	
	spinner := utils.NewLoadingSpinner("Installing PHP...")
	spinner.Start()
	
	output, err := utils.ExecuteCommand("make", "install")
	if err != nil {
		spinner.Stop("❌ Install failed")
		utils.SafeLog(logger, "Install failed: %v", err)
		utils.SafeLog(logger, "Make install output: %s", output)
		return fmt.Errorf("install failed (see log for details)")
	}
	
	spinner.Stop("✓ Install complete")
	utils.SafeLog(logger, "Install completed successfully")
	
	return nil
}