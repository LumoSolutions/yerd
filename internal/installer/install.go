package installer

import (
	"encoding/json"
	"fmt"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/pkg/php"
)

// InstallPHP performs complete PHP installation from source with comprehensive logging and validation.
// version: PHP version to install, uncached: If true, bypasses version cache. Returns error if installation fails.
func InstallPHP(version string, uncached bool) error {
	logger, err := utils.NewLogger(version)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not create log file: %v\n", err)
		logger = nil
	}

	var installSuccess bool
	defer cleanupLogger(logger, &installSuccess)

	if err := setupInstallEnvironment(logger); err != nil {
		return err
	}

	versionInfo, err := fetchAndValidateVersion(version, logger, uncached)
	if err != nil {
		return err
	}

	cfg, err := validateInstallationConfig(version, logger)
	if err != nil {
		return err
	}

	checkSystemPHPConflicts(logger)

	if err := performInstallation(version, versionInfo, logger); err != nil {
		return err
	}

	if err := updateConfiguration(cfg, version, logger); err != nil {
		return err
	}

	installSuccess = true
	return nil
}

// cleanupLogger handles log file cleanup based on installation success status.
// logger: Logger instance, installSuccess: Pointer to success status.
func cleanupLogger(logger *utils.Logger, installSuccess *bool) {
	if logger == nil {
		return
	}

	if *installSuccess {
		logger.DeleteLogFile()
		return
	}

	logPath := logger.Close()
	fmt.Printf("\n📝 Check the detailed installation log:\n")
	fmt.Printf("   %s\n", logPath)
	fmt.Printf("\nTo view the log:\n")
	fmt.Printf("   tail -f %s\n", logPath)
}

// setupInstallEnvironment ensures all required YERD directories exist before installation.
// logger: Logger for operation tracking. Returns error if directory creation fails.
func setupInstallEnvironment(logger *utils.Logger) error {
	utils.SafeLog(logger, "Ensuring YERD directories exist...")

	if err := utils.EnsureDirectories(); err != nil {
		utils.SafeLog(logger, "Directory creation failed: %v", err)
		return fmt.Errorf("failed to create directories: %v", err)
	}

	utils.SafeLog(logger, "Directories created successfully")
	return nil
}

// fetchAndValidateVersion retrieves and validates PHP version information from php.net.
// version: PHP version, logger: Logger instance, uncached: Bypass cache flag. Returns version info or error.
func fetchAndValidateVersion(version string, logger *utils.Logger, uncached bool) (php.VersionInfo, error) {
	utils.SafeLog(logger, "Getting latest PHP version info for: %s", version)
	if uncached {
		utils.SafeLog(logger, "Bypassing cache to get fresh version data")
	}

	var latestVersions map[string]string
	var downloadURLs map[string]string
	var err error

	if uncached {
		latestVersions, downloadURLs, err = versions.GetLatestVersionsFresh()
	} else {
		latestVersions, downloadURLs, err = versions.GetLatestVersions()
	}
	if err != nil {
		utils.SafeLog(logger, "Failed to fetch latest versions from PHP.net: %v", err)
		printVersionFetchError(version)
		return php.VersionInfo{}, fmt.Errorf("cannot install PHP without access to PHP.net API")
	}

	jsonBytes, _ := json.Marshal(latestVersions)
	utils.SafeLog(logger, "Successfully fetched latest versions")
	utils.SafeLog(logger, "latest versions: %s", string(jsonBytes))
	versionInfo, exists := php.GetLatestVersionInfo(version, latestVersions, downloadURLs)

	if !exists {
		utils.SafeLog(logger, "Unsupported PHP version: %s", version)
		return php.VersionInfo{}, fmt.Errorf("unsupported PHP version: %s", version)
	}

	if latestVersion, hasLatest := latestVersions[version]; hasLatest {
		utils.SafeLog(logger, "Installing latest version: %s", latestVersion)
		fmt.Printf("📦 Installing latest PHP %s: %s\n", version, latestVersion)
	}

	utils.SafeLog(logger, "PHP version info found - Source package: %s", versionInfo.SourcePackage)
	return versionInfo, nil
}

// printVersionFetchError displays helpful error message when version fetching fails.
// version: PHP version that failed to fetch.
func printVersionFetchError(version string) {
	fmt.Printf("❌ Could not fetch latest PHP versions from PHP.net\n")
	fmt.Printf("💡 This is required to:\n")
	fmt.Printf("   • Get the latest stable version of PHP %s\n", version)
	fmt.Printf("   • Download the source code from php.net\n")
	fmt.Printf("   • Ensure a secure and up-to-date installation\n\n")
	fmt.Printf("🔍 Troubleshooting:\n")
	fmt.Printf("   • Check your internet connection\n")
	fmt.Printf("   • Verify PHP.net is accessible: curl -I https://www.php.net\n")
	fmt.Printf("   • Try again in a few moments\n")
}

// validateInstallationConfig loads config and checks for existing installations.
// version: PHP version to check, logger: Logger instance. Returns config or error if already installed.
func validateInstallationConfig(version string, logger *utils.Logger) (*config.Config, error) {
	utils.SafeLog(logger, "Loading YERD configuration...")

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.SafeLog(logger, "Failed to load config: %v", err)
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	if _, exists := cfg.InstalledPHP[version]; exists {
		utils.SafeLog(logger, "PHP %s is already installed according to config", version)
		return nil, fmt.Errorf("PHP %s is already installed", version)
	}

	utils.SafeLog(logger, "Configuration loaded successfully")
	return cfg, nil
}

// checkSystemPHPConflicts detects and warns about existing system PHP installations.
// logger: Logger instance for conflict detection logging.
func checkSystemPHPConflicts(logger *utils.Logger) {
	utils.SafeLog(logger, "Checking for existing system PHP...")

	hasSystemPHP, phpType, err := utils.CheckForSystemPHP()
	if err != nil || !hasSystemPHP {
		utils.SafeLog(logger, "No system PHP conflicts detected")
		return
	}

	phpInfo := "Unknown PHP version"
	if info, err := utils.DetectSystemPHPInfo(); err == nil {
		phpInfo = info
	}

	utils.SafeLog(logger, "Existing system PHP detected: %s (%s)", phpInfo, phpType)

	fmt.Printf("⚠️  Warning: Existing PHP installation detected\n")
	fmt.Printf("Found: %s (%s)\n", phpInfo, phpType)
	fmt.Printf("Note: You won't be able to set this as CLI version until you remove the existing PHP\n")
	fmt.Printf("Continuing with installation...\n\n")
}

// performInstallation executes the complete PHP build and installation process.
// version: PHP version, versionInfo: Version details, logger: Logger instance. Returns error if installation fails.
func performInstallation(version string, versionInfo php.VersionInfo, logger *utils.Logger) error {
	installPath := php.GetInstallPath(version)
	binaryPath := php.GetBinaryPath(version)

	utils.SafeLog(logger, "Install path: %s", installPath)
	utils.SafeLog(logger, "Binary path: %s", binaryPath)

	fmt.Printf("Building PHP %s from source...\n", version)

	if err := installFromSource(versionInfo, logger); err != nil {
		utils.SafeLog(logger, "Installation step failed: %v", err)
		return err
	}

	if err := createSymlinks(version, binaryPath, logger); err != nil {
		utils.SafeLog(logger, "Symlink creation failed: %v", err)
		return err
	}

	if err := verifyInstallation(binaryPath, logger); err != nil {
		utils.SafeLog(logger, "Installation verification failed: %v", err)
		return err
	}

	if err := createDefaultPHPIni(version, logger); err != nil {
		utils.SafeLog(logger, "Failed to create default php.ini: %v", err)
		return err
	}

	return nil
}

// updateConfiguration adds the newly installed PHP version to YERD configuration.
// cfg: Config object, version: PHP version, logger: Logger instance. Returns error if config save fails.
func updateConfiguration(cfg *config.Config, version string, logger *utils.Logger) error {
	utils.SafeLog(logger, "Updating YERD configuration...")

	installPath := php.GetInstallPath(version)

	defaultExtensions := getDefaultExtensions()
	cfg.AddInstalledPHPWithExtensions(version, installPath, defaultExtensions)

	if err := cfg.Save(); err != nil {
		utils.SafeLog(logger, "Failed to save config: %v", err)
		return fmt.Errorf("failed to update config: %v", err)
	}

	utils.SafeLog(logger, "Configuration updated successfully")
	utils.SafeLog(logger, "PHP %s installation completed with extensions: %v", version, defaultExtensions)

	return nil
}

// getDefaultExtensions returns the standard set of PHP extensions for new installations.
func getDefaultExtensions() []string {
	return []string{
		"mbstring",
		"bcmath",
		"opcache",
		"curl",
		"openssl",
		"zip",
		"sockets",
		"mysqli",
		"pdo-mysql",
		"gd",
		"jpeg",
		"freetype",
	}
}

// installFromSource handles complete source-based PHP installation with dependency management.
// versionInfo: PHP version and download details, logger: Logger instance. Returns error if build fails.
func installFromSource(versionInfo php.VersionInfo, logger *utils.Logger) error {
	utils.SafeLog(logger, "Starting source installation for PHP %s", versionInfo.Version)
	utils.SafeLog(logger, "Download URL: %s", versionInfo.DownloadURL)

	if err := checkBuildDependencies(logger); err != nil {
		return err
	}

	buildDir := fmt.Sprintf("/tmp/yerd-build-php%s", versionInfo.Version)
	if err := prepareBuildDirectory(buildDir, logger); err != nil {
		return err
	}

	defer func() {
		utils.SafeLog(logger, "Cleaning up build directory: %s", buildDir)
		utils.RemoveDirectory(buildDir)
	}()

	if err := downloadSource(versionInfo, buildDir, logger); err != nil {
		return err
	}

	sourceDir, err := extractSource(versionInfo, buildDir, logger)
	if err != nil {
		return err
	}

	if err := buildAndInstall(versionInfo, sourceDir, logger); err != nil {
		return err
	}

	utils.SafeLog(logger, "Source installation completed successfully")

	fmt.Printf("✓ PHP %s built and installed successfully\n", versionInfo.Version)
	return nil
}
