package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/pkg/php"
)

// prepareBuildDirectory cleans and creates a fresh build directory for PHP source compilation.
// buildDir: Directory path for build operations, logger: Logging instance. Returns error if directory setup fails.
func prepareBuildDirectory(buildDir string, logger *utils.Logger) error {
	if err := os.RemoveAll(buildDir); err != nil && !os.IsNotExist(err) {
		utils.SafeLog(logger, "Failed to clean existing build directory: %v", err)
		return fmt.Errorf("failed to clean build directory: %v", err)
	}

	if err := os.MkdirAll(buildDir, utils.DirPermissions); err != nil {
		utils.SafeLog(logger, "Failed to create build directory: %v", err)
		return fmt.Errorf("failed to create build directory: %v", err)
	}

	utils.SafeLog(logger, "Created build directory: %s", buildDir)
	return nil
}

// downloadSource downloads PHP source archive from official distribution using wget with progress spinner.
// versionInfo: PHP version with download URL, buildDir: Directory for downloaded file, logger: Logging instance. Returns error if download fails.
func downloadSource(versionInfo php.VersionInfo, buildDir string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Downloading PHP source from: %s", versionInfo.DownloadURL)

	filename := versionInfo.SourcePackage + ".tar.gz"
	filePath := filepath.Join(buildDir, filename)

	spinnerMessage := fmt.Sprintf("Downloading %s...", filename)
	spinner := utils.NewLoadingSpinner(spinnerMessage)
	spinner.Start()

	output, err := utils.ExecuteCommand("wget", "-O", filePath, versionInfo.DownloadURL)
	if err != nil {
		spinner.Stop("❌ Download failed")
		utils.SafeLog(logger, "Download failed: %v", err)
		utils.SafeLog(logger, "wget output: %s", output)
		return fmt.Errorf("failed to download PHP source: %v", err)
	}

	spinner.Stop("✓ Download complete")
	utils.SafeLog(logger, "Downloaded to: %s", filePath)

	return nil
}

// extractSource extracts PHP source archive to build directory and returns extracted source path.
// versionInfo: PHP version with package info, buildDir: Build directory path, logger: Logging instance. Returns source directory or error.
func extractSource(versionInfo php.VersionInfo, buildDir string, logger *utils.Logger) (string, error) {
	utils.SafeLog(logger, "Extracting PHP source archive")

	filename := versionInfo.SourcePackage + ".tar.gz"
	filePath := filepath.Join(buildDir, filename)

	spinner := utils.NewLoadingSpinner("Extracting source code...")
	spinner.Start()

	output, err := utils.ExecuteCommand("tar", "-xzf", filePath, "-C", buildDir)
	if err != nil {
		spinner.Stop("❌ Extraction failed")
		utils.SafeLog(logger, "Extraction failed: %v", err)
		utils.SafeLog(logger, "tar output: %s", output)
		return "", fmt.Errorf("failed to extract PHP source: %v", err)
	}

	sourceDir := filepath.Join(buildDir, versionInfo.SourcePackage)
	spinner.Stop("✓ Source extracted")

	utils.SafeLog(logger, "Extracted to: %s", sourceDir)

	return sourceDir, nil
}
