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

// SetCLIVersion creates symlinks to set a PHP version as the system default CLI.
// version: PHP version string to set as CLI. Returns error if version not installed or symlink fails.
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

// RemovePHP completely removes a PHP version including symlinks, CLI settings, and source files.
// version: PHP version string to remove. Returns error if version not installed or removal fails.
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

// removeSymlinks removes global and version-specific symlinks for a PHP installation.
// version: PHP version string whose symlinks should be removed.
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

// handleCLIRemoval manages CLI symlink removal and provides user guidance for alternatives.
// cfg: Configuration object, version: PHP version being removed as CLI.
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

// removeSourceInstallation deletes the PHP source installation directory.
// cfg: Configuration object, version: PHP version whose installation should be removed.
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

// ListInstalledVersions returns a slice of all currently installed PHP version strings.
// Returns version slice or error if config loading fails.
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

// GetCurrentCLIVersion returns the PHP version currently set as CLI default.
// Returns version string (empty if none set) or error if config loading fails.
func GetCurrentCLIVersion() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf(configLoadError, err)
	}

	return cfg.CurrentCLI, nil
}

// getRemainingVersions returns all installed PHP versions except the specified one.
// cfg: Configuration object, excludeVersion: Version to exclude. Returns filtered version slice.
func getRemainingVersions(cfg *config.Config, excludeVersion string) []string {
	var versions []string
	for version := range cfg.InstalledPHP {
		if version != excludeVersion {
			versions = append(versions, version)
		}
	}
	return versions
}
