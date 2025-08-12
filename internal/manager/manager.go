package manager

import (
	"fmt"
	"os"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/pkg/php"
)

const (
	configLoadError = "failed to load config: %v"
	globalPHPPath   = "/usr/local/bin/php"
)

func SetCLIVersion(version string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf(configLoadError, err)
	}
	
	if _, exists := cfg.InstalledPHP[version]; !exists {
		return fmt.Errorf("PHP %s is not installed", version)
	}
	
	binaryPath := php.GetBinaryPath(version)
	if !utils.FileExists(binaryPath) {
		return fmt.Errorf("PHP binary not found at %s", binaryPath)
	}
	
	// Remove any existing file at CLI path (symlink or regular file) before creating new one
	if _, err := os.Lstat(globalPHPPath); err == nil {
		if err := os.Remove(globalPHPPath); err != nil {
			return fmt.Errorf("failed to remove existing CLI file: %v", err)
		}
	}
	
	if err := utils.CreateSymlink(binaryPath, globalPHPPath); err != nil {
		return fmt.Errorf("failed to create CLI symlink: %v", err)
	}
	
	cfg.SetCurrentCLI(version)
	if err := cfg.Save(); err != nil {
		return fmt.Errorf(configLoadError, err)
	}
	
	return nil
}

func RemovePHP(version string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf(configLoadError, err)
	}
	
	if _, exists := cfg.InstalledPHP[version]; !exists {
		return fmt.Errorf("PHP %s is not installed", version)
	}
	
	removeSymlinks(version)
	
	if cfg.CurrentCLI == version {
		handleCLIRemoval(cfg, version)
	}
	
	removeSourceInstallation(cfg, version)
	
	cfg.RemoveInstalledPHP(version)
	if err := cfg.Save(); err != nil {
		return fmt.Errorf(configLoadError, err)
	}
	
	return nil
}

func removeSymlinks(version string) {
	fmt.Printf("Removing symlinks...\n")
	
	binaryPath := php.GetBinaryPath(version)
	globalBinaryPath := globalPHPPath + version
	
	if err := utils.RemoveSymlink(globalBinaryPath); err != nil {
		fmt.Printf("Warning: failed to remove global binary symlink: %v\n", err)
	}
	
	if err := utils.RemoveSymlink(binaryPath); err != nil {
		fmt.Printf("Warning: failed to remove binary symlink: %v\n", err)
	}
}

func handleCLIRemoval(cfg *config.Config, version string) {
	fmt.Printf("⚠️  Removing current CLI version (PHP %s)\n", version)
	
	if err := utils.RemoveSymlink(globalPHPPath); err != nil {
		fmt.Printf("Warning: failed to remove CLI symlink: %v\n", err)
	} else {
		fmt.Printf("✓ Removed CLI version symlink (%s)\n", globalPHPPath)
	}
	
	remainingVersions := getRemainingVersions(cfg, version)
	if len(remainingVersions) > 0 {
		fmt.Printf("💡 Suggestion: Set a new CLI version with: sudo yerd php cli %s\n", remainingVersions[0])
	} else {
		fmt.Printf("ℹ️  No other PHP versions installed. Install one with: sudo yerd php add 8.4\n")
	}
}

func removeSourceInstallation(cfg *config.Config, version string) {
	fmt.Printf("Removing source installation...\n")
	
	phpInfo, exists := cfg.InstalledPHP[version]
	if !exists || phpInfo.InstallPath == "" {
		return
	}
	
	if err := utils.RemoveDirectory(phpInfo.InstallPath); err != nil {
		fmt.Printf("Warning: failed to remove installation directory: %v\n", err)
	} else {
		fmt.Printf("✓ Removed PHP %s installation directory\n", version)
	}
}


func ListInstalledVersions() ([]string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf(configLoadError, err)
	}
	
	var versions []string
	for version := range cfg.InstalledPHP {
		versions = append(versions, version)
	}
	
	return versions, nil
}

func GetCurrentCLIVersion() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf(configLoadError, err)
	}
	
	return cfg.CurrentCLI, nil
}

func getRemainingVersions(cfg *config.Config, excludeVersion string) []string {
	var versions []string
	for version := range cfg.InstalledPHP {
		if version != excludeVersion {
			versions = append(versions, version)
		}
	}
	return versions
}

