package php

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/pkg/php"
)

// Common error message templates
const (
	ErrInvalidPHPVersion     = "invalid PHP version: %s"
	ErrFailedToLoadConfig    = "failed to load config: %v"
	ErrPHPNotInstalled       = "PHP %s is not installed. Use 'yerd php add %s' first"
	ErrNoExtensionInfo       = "no extension information found for PHP %s"
	ErrFailedToCreateBuilder = "failed to create builder: %v"
	ErrRebuildFailed         = "rebuild failed"
)

// Common operation descriptions
const (
	OpExtensionManagement = "Extension management"
	OpPHPRebuild         = "PHP rebuild"
	OpPHPInstall         = "PHP installation"
	OpPHPRemoval         = "PHP removal"
)

// ValidatePHPVersionAndConfig validates PHP version and loads config with standardized error messages.
// version: PHP version string. Returns formatted version, loaded config, and error if validation fails.
func ValidatePHPVersionAndConfig(version string) (string, *config.Config, error) {
	phpVersion := php.FormatVersion(version)

	if !php.IsValidVersion(phpVersion) {
		return "", nil, fmt.Errorf(ErrInvalidPHPVersion, phpVersion)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return phpVersion, nil, fmt.Errorf(ErrFailedToLoadConfig, err)
	}

	return phpVersion, cfg, nil
}

// ValidateInstalledPHP validates that a PHP version is installed with standardized error message.
// cfg: Config object, version: PHP version string. Returns error if not installed.
func ValidateInstalledPHP(cfg *config.Config, version string) error {
	if _, exists := cfg.InstalledPHP[version]; !exists {
		return fmt.Errorf(ErrPHPNotInstalled, version, version)
	}
	return nil
}

// ValidatePHPVersionConfigAndInstallation combines all common PHP validation steps.
// version: PHP version string. Returns formatted version, loaded config, and error if any validation fails.
func ValidatePHPVersionConfigAndInstallation(version string) (string, *config.Config, error) {
	phpVersion, cfg, err := ValidatePHPVersionAndConfig(version)
	if err != nil {
		return phpVersion, cfg, err
	}

	if err := ValidateInstalledPHP(cfg, phpVersion); err != nil {
		return phpVersion, cfg, err
	}

	return phpVersion, cfg, nil
}