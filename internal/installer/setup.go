package installer

import (
	"fmt"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
)

func createSymlinks(version, binaryPath string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Creating symlinks for PHP %s", version)
	utils.SafeLog(logger, "Target binary path: %s", binaryPath)
	
	spinner := utils.NewLoadingSpinner(fmt.Sprintf("Locating PHP %s binary...", version))
	spinner.Start()
	
	systemBinaryPath, err := findInstalledPHPBinary(version, logger)
	if err != nil {
		spinner.Stop("❌ Binary not found")
		utils.SafeLog(logger, "Failed to locate installed PHP binary: %v", err)
		return fmt.Errorf("PHP binary not found after installation (check log for details)")
	}
	
	spinner.Stop(fmt.Sprintf("✓ Found PHP binary at: %s", systemBinaryPath))
	utils.SafeLog(logger, "Located system PHP binary at: %s", systemBinaryPath)
	
	symlinkSpinner := utils.NewLoadingSpinner("Creating symlinks...")
	symlinkSpinner.Start()
	
	utils.SafeLog(logger, "Creating YERD binary symlink: %s -> %s", binaryPath, systemBinaryPath)
	
	if err := utils.CreateSymlink(systemBinaryPath, binaryPath); err != nil {
		symlinkSpinner.Stop("❌ Symlink creation failed")
		utils.SafeLog(logger, "Failed to create binary symlink: %v", err)
		return fmt.Errorf("symlink creation failed (check permissions)")
	}
	
	globalBinaryPath := utils.SystemBinDir + "/php" + version
	utils.SafeLog(logger, "Creating global binary symlink: %s -> %s", globalBinaryPath, binaryPath)
	
	if err := utils.CreateSymlink(binaryPath, globalBinaryPath); err != nil {
		symlinkSpinner.Stop("❌ Global symlink creation failed")
		utils.SafeLog(logger, "Failed to create global symlink: %v", err)
		return fmt.Errorf("global symlink creation failed (check sudo permissions)")
	}
	
	symlinkSpinner.Stop("✓ Symlinks created successfully")
	utils.SafeLog(logger, "All symlinks created successfully")
	
	return nil
}

func verifyInstallation(binaryPath string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Verifying PHP installation at: %s", binaryPath)
	
	if err := utils.VerifyPHPInstallation(binaryPath); err != nil {
		utils.SafeLog(logger, "PHP installation verification failed: %v", err)
		return err
	}
	
	utils.SafeLog(logger, "PHP installation verification successful")
	
	fmt.Printf("✓ PHP installation verified\n")
	return nil
}

func createDefaultPHPIni(version string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Creating default php.ini for PHP %s", version)
	
	err := utils.CreatePHPIniForVersion(version)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.SafeLog(logger, "php.ini already exists for PHP %s", version)
			fmt.Printf("✓ php.ini already exists\n")
			return nil
		}
		
		utils.SafeLog(logger, "Failed to create php.ini: %v", err)
		return err
	}
	
	utils.SafeLog(logger, "Default php.ini created successfully for PHP %s", version)
	
	fmt.Printf("✓ Default php.ini created\n")
	return nil
}