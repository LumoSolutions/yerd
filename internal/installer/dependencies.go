package installer

import (
	"fmt"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
)

func checkBuildDependencies(logger *utils.Logger) error {
	utils.SafeLog(logger, "Checking build dependencies...")
	
	distro, err := utils.DetectDistro()
	if err != nil {
		utils.SafeLog(logger, "Failed to detect distribution: %v", err)
		return fmt.Errorf("failed to detect Linux distribution: %v", err)
	}
	
	utils.SafeLog(logger, "Detected distribution: %s (%s)", distro.Name, distro.PackageManager)
	
	fmt.Printf("📋 Detected: %s\n", distro.Name)
	
	dependencyKeys := []string{
		"gcc", "make", "autoconf", "pkgconf", "bison", "re2c",
		"libxml2", "curl", "openssl", "zlib", "libpng", "libjpeg-turbo", "freetype2",
	}
	
	var missing []string
	
	for _, depKey := range dependencyKeys {
		packageName := utils.GetPackageForDistro(depKey, distro)
		
		utils.SafeLog(logger, "Checking dependency: %s -> %s", depKey, packageName)
		
		var isInstalled bool
		var checkErr error
		
		switch distro.PackageManager {
		case "pacman":
			_, checkErr = utils.ExecuteCommand(distro.QueryCmd[0], distro.QueryCmd[1], packageName)
			isInstalled = checkErr == nil
		case "apt":
			output, checkErr := utils.ExecuteCommand(distro.QueryCmd[0], distro.QueryCmd[1], packageName)
			isInstalled = checkErr == nil && strings.Contains(output, packageName)
		case "dnf", "yum", "zypper":
			_, checkErr = utils.ExecuteCommand(distro.QueryCmd[0], distro.QueryCmd[1], packageName)
			isInstalled = checkErr == nil
		default:
			utils.SafeLog(logger, "Unknown package manager: %s", distro.PackageManager)
			return fmt.Errorf("unsupported package manager: %s", distro.PackageManager)
		}
		
		if !isInstalled {
			missing = append(missing, packageName)
			utils.SafeLog(logger, "Missing dependency: %s (%s)", depKey, packageName)
		} else {
			utils.SafeLog(logger, "Found dependency: %s (%s)", depKey, packageName)
		}
	}
	
	if len(missing) > 0 {
		fmt.Printf("⚠️  Installing build dependencies...\n")
		utils.SafeLog(logger, "Installing missing dependencies: %v", missing)
		
		if distro.PackageManager == "apt" {
			fmt.Printf("📦 Updating package cache...\n")
			utils.SafeLog(logger, "Updating apt package cache...")
			output, err := utils.ExecuteCommandWithLogging(logger, "apt-get", "update")
			if err != nil {
				utils.SafeLog(logger, "Failed to update package cache: %v", err)
				utils.SafeLog(logger, "apt-get update output: %s", output)
				return fmt.Errorf("failed to update package cache: %v", err)
			}
		}
		
		installCmd := append(distro.InstallCmd, missing...)
		output, err := utils.ExecuteCommandWithLogging(logger, installCmd[0], installCmd[1:]...)
		if err != nil {
			utils.SafeLog(logger, "Failed to install dependencies: %v", err)
			utils.SafeLog(logger, "%s output: %s", distro.PackageManager, output)
			return fmt.Errorf("failed to install build dependencies using %s: %v", distro.PackageManager, err)
		}
		
		fmt.Printf("✓ Build dependencies installed\n")
	} else {
		fmt.Printf("✓ All build dependencies satisfied\n")
	}
	
	return nil
}