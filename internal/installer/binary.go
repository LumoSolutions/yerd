package installer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
)

const binaryVerificationError = "Binary verification failed for %s: %v"

// findInstalledPHPBinary locates the PHP binary after source installation with fallback paths.
// version: PHP version to locate, logger: Logging instance. Returns binary path or error if not found.
func findInstalledPHPBinary(version string, logger *utils.Logger) (string, error) {
	utils.SafeLog(logger, "Locating PHP %s binary after source build...", version)

	expectedPath := utils.YerdPHPDir + "/php" + version + "/bin/php"

	utils.SafeLog(logger, "Checking expected source build path: %s", expectedPath)

	if utils.FileExists(expectedPath) {
		utils.SafeLog(logger, "File exists at %s, verifying...", expectedPath)
		if err := verifyPHPBinary(expectedPath, version, logger); err == nil {
			utils.SafeLog(logger, "Found and verified PHP binary at: %s", expectedPath)
			return expectedPath, nil
		} else {
			utils.SafeLog(logger, binaryVerificationError, expectedPath, err)
		}
	} else {
		utils.SafeLog(logger, "Binary not found at expected path: %s", expectedPath)
	}

	alternativePaths := []string{
		utils.YerdPHPDir + "/php" + version + "/bin/php-cli",
		utils.YerdPHPDir + "/php" + version + "/sbin/php-fpm",
		utils.SystemBinDir + "/php" + version,
		utils.SystemBinDir + "/php",
	}

	utils.SafeLog(logger, "Checking alternative paths: %v", alternativePaths)

	for _, path := range alternativePaths {
		utils.SafeLog(logger, "Checking path: %s", path)
		if utils.FileExists(path) {
			utils.SafeLog(logger, "File exists at %s, verifying...", path)
			if err := verifyPHPBinary(path, version, logger); err == nil {
				utils.SafeLog(logger, "Found and verified PHP binary at: %s", path)
				return path, nil
			} else {
				utils.SafeLog(logger, binaryVerificationError, path, err)
			}
		}
	}

	installDir := utils.YerdPHPDir + "/php" + version
	utils.SafeLog(logger, "Searching installation directory for PHP binary: %s", installDir)

	if foundPath, err := searchForPHPInInstallDir(installDir, version, logger); err == nil && foundPath != "" {
		utils.SafeLog(logger, "Found PHP binary in install directory: %s", foundPath)
		return foundPath, nil
	} else {
		utils.SafeLog(logger, "Search in install directory failed: %v", err)
	}

	utils.SafeLog(logger, "PHP %s binary not found after source installation", version)
	fmt.Printf("🔍 Debug: Searching for PHP %s binary in installation directory...\n", version)
	showInstalledSourceFiles(version, installDir, logger)

	return "", fmt.Errorf("PHP %s binary not found after source installation", version)
}

// searchForPHPInInstallDir recursively searches installation directory for PHP executable binary.
// installDir: Directory to search, version: Expected PHP version, logger: Logging instance. Returns binary path or error.
func searchForPHPInInstallDir(installDir, version string, logger *utils.Logger) (string, error) {
	utils.SafeLog(logger, "Searching for PHP binary in: %s", installDir)

	output, err := utils.ExecuteCommand("find", installDir, "-name", "php", "-type", "f", "-executable")
	if err != nil {
		utils.SafeLog(logger, "Find command failed: %v", err)
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		utils.SafeLog(logger, "Found potential PHP binary: %s", line)

		if err := verifyPHPBinary(line, version, logger); err == nil {
			utils.SafeLog(logger, "Verified PHP binary: %s", line)
			return line, nil
		} else {
			utils.SafeLog(logger, binaryVerificationError, line, err)
		}
	}

	return "", fmt.Errorf("no valid PHP binary found in %s", installDir)
}

// verifyPHPBinary validates that a PHP binary exists and reports the expected version.
// path: Binary path to verify, expectedVersion: Required PHP version, logger: Logging instance. Returns error if invalid.
func verifyPHPBinary(path, expectedVersion string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Verifying PHP binary at %s for version %s", path, expectedVersion)

	output, err := utils.ExecuteCommand(path, "-v")
	if err != nil {
		utils.SafeLog(logger, "Failed to execute PHP binary %s: %v", path, err)
		return fmt.Errorf("failed to execute PHP binary: %v", err)
	}

	utils.SafeLog(logger, "PHP binary output: %s", strings.TrimSpace(output))

	if !strings.Contains(output, "PHP "+expectedVersion) {
		utils.SafeLog(logger, "PHP version mismatch: expected %s, got: %s", expectedVersion, output)
		return fmt.Errorf("PHP version mismatch: expected %s, got: %s", expectedVersion, output)
	}

	utils.SafeLog(logger, "PHP binary verification successful")

	return nil
}

// searchForPHPInDir searches a directory for PHP binaries using common naming patterns.
// dir: Directory to search, version: PHP version to find, logger: Logging instance. Returns binary path or error.
func searchForPHPInDir(dir, version string, logger *utils.Logger) (string, error) {
	patterns := []string{
		"php" + version,
		"php-" + version,
		"php" + version + "-cli",
	}

	utils.SafeLog(logger, "Searching directory %s for patterns: %v", dir, patterns)

	for _, pattern := range patterns {
		path := filepath.Join(dir, pattern)
		utils.SafeLog(logger, "Checking: %s", path)
		if utils.FileExists(path) {
			utils.SafeLog(logger, "File exists, verifying: %s", path)
			if err := verifyPHPBinary(path, version, logger); err == nil {
				utils.SafeLog(logger, "Found and verified PHP binary: %s", path)
				return path, nil
			} else {
				utils.SafeLog(logger, "Verification failed for %s: %v", path, err)
			}
		}
	}

	utils.SafeLog(logger, "No matching PHP binary found in %s", dir)

	return "", fmt.Errorf("no matching PHP binary found in %s", dir)
}

// showInstalledSourceFiles displays debug information about files in installation and system directories.
// version: PHP version being searched, installDir: Installation directory path, logger: Logging instance.
func showInstalledSourceFiles(version string, installDir string, logger *utils.Logger) {
	fmt.Printf("Searching for PHP %s files in installation directory:\n", version)
	utils.SafeLog(logger, "Running debug search for installed PHP files in: %s", installDir)

	searchInstallationDirectory(installDir, logger)
	searchSystemDirectories(version, logger)
}

// searchInstallationDirectory searches the main installation directory for executables and PHP files.
// installDir: Installation directory path, logger: Logging instance.
func searchInstallationDirectory(installDir string, logger *utils.Logger) {
	if !utils.FileExists(installDir) {
		fmt.Printf("❌ Installation directory does not exist: %s\n", installDir)
		utils.SafeLog(logger, "Installation directory does not exist: %s", installDir)
		return
	}

	fmt.Printf("📁 Installation directory: %s\n", installDir)
	utils.SafeLog(logger, "Searching in installation directory: %s", installDir)

	searchExecutableFiles(installDir, logger)
	searchPHPFiles(installDir, logger)
}

// searchExecutableFiles finds and displays all executable files in a directory.
// dir: Directory to search, logger: Logging instance.
func searchExecutableFiles(dir string, logger *utils.Logger) {
	output, err := utils.ExecuteCommand("find", dir, "-type", "f", "-executable")
	if err != nil || output == "" {
		fmt.Printf("   No executable files found\n")
		utils.SafeLog(logger, "No executable files found: %v", err)
		return
	}

	fmt.Printf("   Executable files found:\n")
	printFileList(output, "   - ", logger, "Found executable")
}

// searchPHPFiles finds and displays all PHP-related files in a directory.
// dir: Directory to search, logger: Logging instance.
func searchPHPFiles(dir string, logger *utils.Logger) {
	output, err := utils.ExecuteCommand("find", dir, "-name", "*php*", "-type", "f")
	if err != nil || output == "" {
		utils.SafeLog(logger, "No PHP files found: %v", err)
		return
	}

	fmt.Printf("   PHP-related files:\n")
	printFileList(output, "   - ", logger, "Found PHP file")
}

// searchSystemDirectories searches standard system directories for PHP binaries.
// version: PHP version to search for, logger: Logging instance.
func searchSystemDirectories(version string, logger *utils.Logger) {
	fmt.Printf("\n🔍 Checking system directories for PHP %s:\n", version)
	searchDirs := []string{"/usr/local/bin", "/usr/local/sbin"}

	for _, dir := range searchDirs {
		searchSystemDirectory(dir, logger)
	}
}

// searchSystemDirectory searches a specific system directory for PHP-related files.
// dir: System directory path to search, logger: Logging instance.
func searchSystemDirectory(dir string, logger *utils.Logger) {
	utils.SafeLog(logger, "Searching for PHP files in: %s", dir)
	output, err := utils.ExecuteCommand("find", dir, "-name", "*php*", "-type", "f", "2>/dev/null")
	if err != nil || output == "" {
		utils.SafeLog(logger, "Find command failed for %s: %v", dir, err)
		return
	}

	printFileList(output, fmt.Sprintf("  Found in %s: ", dir), logger, "Found in system dir")
}

// printFileList formats and displays a list of files with consistent prefix and logging.
// output: Newline-separated file list, prefix: Display prefix, logger: Logging instance, logMessage: Log entry prefix.
func printFileList(output, prefix string, logger *utils.Logger, logMessage string) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			fmt.Printf("%s%s\n", prefix, line)
			utils.SafeLog(logger, "%s: %s", logMessage, line)
		}
	}
}
