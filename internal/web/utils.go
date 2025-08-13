package web

import (
	"fmt"
	"path/filepath"

	"github.com/LumoSolutions/yerd/internal/utils"
)

// VerifyServiceInstallation checks if a web service is properly installed
func VerifyServiceInstallation(serviceName string) error {
	_, exists := GetServiceConfig(serviceName)
	if !exists {
		return fmt.Errorf("unsupported service: %s", serviceName)
	}

	binaryPath := GetServiceBinaryPath(serviceName)
	if !utils.FileExists(binaryPath) {
		return fmt.Errorf("service binary not found at %s", binaryPath)
	}

	_, err := utils.ExecuteCommand(binaryPath, "-v")
	if err != nil {
		return fmt.Errorf("service binary is not executable or corrupted: %w", err)
	}

	return nil
}

// GetServiceVersion returns the version of an installed service
func GetServiceVersion(serviceName string) (string, error) {
	binaryPath := GetServiceBinaryPath(serviceName)
	if !utils.FileExists(binaryPath) {
		return "", fmt.Errorf("service not installed")
	}

	output, err := utils.ExecuteCommand(binaryPath, "-v")
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}

	return output, nil
}

// IsServiceInstalled checks if a service is already installed
func IsServiceInstalled(serviceName string) bool {
	binaryPath := GetServiceBinaryPath(serviceName)
	return utils.FileExists(binaryPath)
}

// RemoveService removes an installed web service
func RemoveService(serviceName string) error {
	if !IsServiceInstalled(serviceName) {
		return fmt.Errorf("service %s is not installed", serviceName)
	}

	installPath := GetServiceInstallPath(serviceName)
	return utils.RemoveDirectory(installPath)
}

// CreateServiceSymlink creates a symlink for a service binary in system bin
func CreateServiceSymlink(serviceName string) error {
	binaryPath := GetServiceBinaryPath(serviceName)
	if !utils.FileExists(binaryPath) {
		return fmt.Errorf("service binary not found at %s", binaryPath)
	}

	symlinkPath := filepath.Join(utils.SystemBinDir, serviceName)
	return utils.CreateSymlink(binaryPath, symlinkPath)
}

// RemoveServiceSymlink removes a service symlink from system bin
func RemoveServiceSymlink(serviceName string) error {
	symlinkPath := filepath.Join(utils.SystemBinDir, serviceName)
	return utils.RemoveSymlink(symlinkPath)
}
