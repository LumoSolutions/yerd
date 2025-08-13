package manager

import (
	"fmt"
	"os"

	"github.com/LumoSolutions/yerd/internal/utils"
)

const (
	globalComposerPath = "/usr/local/bin/composer"
	composerPharName   = "composer.phar"
)

// InstallComposer downloads and installs Composer in YERD directory structure with global symlink.
// Returns error if download fails, installation fails, or symlink creation fails.
func InstallComposer() error {
	if err := utils.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to create YERD directories: %v", err)
	}

	composerPath := getComposerInstallPath()

	if err := downloadComposer(composerPath); err != nil {
		return fmt.Errorf("failed to download composer: %v", err)
	}

	if err := makeComposerExecutable(composerPath); err != nil {
		return fmt.Errorf("failed to make composer executable: %v", err)
	}

	if err := createComposerSymlink(composerPath); err != nil {
		return fmt.Errorf("failed to create composer symlink: %v", err)
	}

	return nil
}

// RemoveComposer removes Composer installation and global symlink from the system.
// Returns error if removal fails for composer.phar file or global symlink.
func RemoveComposer() error {
	composerPath := getComposerInstallPath()

	if err := removeComposerFile(composerPath); err != nil {
		return fmt.Errorf("failed to remove composer file: %v", err)
	}

	if err := removeComposerGlobalSymlink(); err != nil {
		return fmt.Errorf("failed to remove composer symlink: %v", err)
	}

	return nil
}

// getComposerInstallPath returns the full path where composer.phar should be installed.
// Returns path string combining YERD bin directory and composer phar filename.
func getComposerInstallPath() string {
	return utils.YerdBinDir + "/" + composerPharName
}

// downloadComposer downloads the latest composer.phar from the official installer.
// composerPath: Target path where composer.phar should be saved. Returns error if download fails.
func downloadComposer(composerPath string) error {
	spinner := utils.NewLoadingSpinner("Downloading Composer...")
	spinner.Start()
	defer spinner.Stop("")

	downloadURL := "https://getcomposer.org/download/latest-stable/composer.phar"

	if err := downloadWithAvailableTool(composerPath, downloadURL); err != nil {
		spinner.Stop("❌ Download failed")
		return fmt.Errorf("download failed: %v", err)
	}

	if !utils.FileExists(composerPath) {
		spinner.Stop("❌ Download failed")
		return fmt.Errorf("composer.phar not found after download")
	}

	spinner.Stop("✓ Composer downloaded")
	return nil
}

// downloadWithAvailableTool attempts to download using curl or wget as fallback.
// composerPath: Target file path, downloadURL: Source URL. Returns error if both tools fail.
func downloadWithAvailableTool(composerPath, downloadURL string) error {
	if isCommandAvailable("curl") {
		return downloadWithCurl(composerPath, downloadURL)
	}

	if isCommandAvailable("wget") {
		return downloadWithWget(composerPath, downloadURL)
	}

	return fmt.Errorf("neither curl nor wget is available on this system")
}

// isCommandAvailable checks if a command is available in the system PATH.
// command: Command name to check. Returns true if command is available.
func isCommandAvailable(command string) bool {
	output, err := utils.ExecuteCommand("which", command)
	return err == nil && output != ""
}

// downloadWithCurl downloads a file using curl with appropriate flags.
// composerPath: Target file path, downloadURL: Source URL. Returns error if download fails.
func downloadWithCurl(composerPath, downloadURL string) error {
	output, err := utils.ExecuteCommand("curl", "-sS", "-L", "-o", composerPath, downloadURL)
	if err != nil {
		return fmt.Errorf("curl failed: %v, output: %s", err, output)
	}
	return nil
}

// downloadWithWget downloads a file using wget with appropriate flags.
// composerPath: Target file path, downloadURL: Source URL. Returns error if download fails.
func downloadWithWget(composerPath, downloadURL string) error {
	output, err := utils.ExecuteCommand("wget", "-q", "-O", composerPath, downloadURL)
	if err != nil {
		return fmt.Errorf("wget failed: %v, output: %s", err, output)
	}
	return nil
}

// makeComposerExecutable sets executable permissions on the composer.phar file.
// composerPath: Path to composer.phar file. Returns error if chmod fails.
func makeComposerExecutable(composerPath string) error {
	spinner := utils.NewLoadingSpinner("Setting permissions...")
	spinner.Start()
	defer spinner.Stop("")

	if err := os.Chmod(composerPath, 0755); err != nil {
		spinner.Stop("❌ Permission setting failed")
		return fmt.Errorf("failed to chmod composer: %v", err)
	}

	spinner.Stop("✓ Permissions set")
	return nil
}

// createComposerSymlink creates a global symlink for composer command accessibility.
// composerPath: Path to composer.phar file. Returns error if symlink creation fails.
func createComposerSymlink(composerPath string) error {
	spinner := utils.NewLoadingSpinner("Creating global symlink...")
	spinner.Start()
	defer spinner.Stop("")

	if _, err := os.Lstat(globalComposerPath); err == nil {
		if err := os.Remove(globalComposerPath); err != nil {
			spinner.Stop("❌ Symlink creation failed")
			return fmt.Errorf("failed to remove existing composer symlink: %v", err)
		}
	}

	if err := utils.CreateSymlink(composerPath, globalComposerPath); err != nil {
		spinner.Stop("❌ Symlink creation failed")
		return fmt.Errorf("failed to create composer symlink: %v", err)
	}

	spinner.Stop("✓ Global symlink created")
	return nil
}

// removeComposerFile removes the composer.phar file if it exists.
// composerPath: Path to composer.phar file. Returns error if removal fails.
func removeComposerFile(composerPath string) error {
	spinner := utils.NewLoadingSpinner("Removing Composer file...")
	spinner.Start()
	defer spinner.Stop("")

	if !utils.FileExists(composerPath) {
		spinner.Stop("⚠️  Composer file not found")
		return nil
	}

	if err := os.Remove(composerPath); err != nil {
		spinner.Stop("❌ File removal failed")
		return fmt.Errorf("failed to remove composer.phar: %v", err)
	}

	spinner.Stop("✓ Composer file removed")
	return nil
}

// removeComposerGlobalSymlink removes the global composer symlink if it exists.
// Returns error if symlink removal fails.
func removeComposerGlobalSymlink() error {
	spinner := utils.NewLoadingSpinner("Removing global symlink...")
	spinner.Start()
	defer spinner.Stop("")

	if _, err := os.Lstat(globalComposerPath); err != nil {
		spinner.Stop("⚠️  Global symlink not found")
		return nil
	}

	if err := utils.RemoveSymlink(globalComposerPath); err != nil {
		spinner.Stop("❌ Symlink removal failed")
		return fmt.Errorf("failed to remove composer symlink: %v", err)
	}

	spinner.Stop("✓ Global symlink removed")
	return nil
}
