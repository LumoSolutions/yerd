package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const phpIniFile = "/php.ini"

// Helper functions for common operations

// validatePHPVersion checks if PHP version is not empty and returns an error if it is.
// version: PHP version string. Returns error if empty.
func validatePHPVersion(version string) error {
	if version == "" {
		return fmt.Errorf(ErrEmptyPHPVersion)
	}
	return nil
}

// getFPMPoolConfigDir returns the FPM pool configuration directory path for a PHP version.
// version: PHP version string. Returns full directory path.
func getFPMPoolConfigDir(version string) string {
	configDir := filepath.Join(YerdEtcDir, "php"+version)
	return filepath.Join(configDir, FPMPoolDir)
}

// getFPMPoolConfigFile returns the FPM pool configuration file path for a PHP version.
// version: PHP version string. Returns full file path.
func getFPMPoolConfigFile(version string) string {
	return filepath.Join(getFPMPoolConfigDir(version), FPMPoolConfig)
}

// GetFPMPoolConfigFile returns the FPM pool configuration file path for a PHP version (public version).
// version: PHP version string. Returns full file path.
func GetFPMPoolConfigFile(version string) string {
	return getFPMPoolConfigFile(version)
}

// getFPMConfigDir returns the FPM main configuration directory path for a PHP version.
// version: PHP version string. Returns directory path.
func getFPMConfigDir(version string) string {
	return filepath.Join(YerdEtcDir, "php"+version)
}

// getFPMMainConfigFile returns the FPM main configuration file path for a PHP version.
// version: PHP version string. Returns full file path.
func getFPMMainConfigFile(version string) string {
	return filepath.Join(getFPMConfigDir(version), "php-fpm.conf")
}

// GetFPMMainConfigFile returns the FPM main configuration file path for a PHP version (public version).
// version: PHP version string. Returns full file path.
func GetFPMMainConfigFile(version string) string {
	return getFPMMainConfigFile(version)
}

// getSystemdServiceName returns the systemd service name for a PHP version.
// version: PHP version string. Returns service name.
func getSystemdServiceName(version string) string {
	return fmt.Sprintf("yerd-php%s-fpm.service", version)
}

// getSystemdServicePath returns the systemd service file path for a PHP version.
// version: PHP version string. Returns full service file path.
func getSystemdServicePath(version string) string {
	return filepath.Join(SystemdDir, getSystemdServiceName(version))
}

// GetSystemdServicePath returns the systemd service file path for a PHP version (public version).
// version: PHP version string. Returns full service file path.
func GetSystemdServicePath(version string) string {
	return getSystemdServicePath(version)
}

type PHPVersionDetails struct {
	Version    string
	BinaryPath string
	IniPath    string
	IsWorking  bool
	ErrorMsg   string
}

// CheckForSystemPHP detects if PHP is installed outside of YERD management.
// Returns hasSystemPHP boolean, PHP type description, and error if detection fails.
func CheckForSystemPHP() (bool, string, error) {
	globalPHPPath := filepath.Join(SystemBinDir, "php")

	if !FileExists(globalPHPPath) {
		return false, "", nil
	}

	info, err := os.Lstat(globalPHPPath)
	if err != nil {
		return false, "", fmt.Errorf("failed to check php binary: %v", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(globalPHPPath)
		if err != nil {
			return false, "", fmt.Errorf("failed to read symlink target: %v", err)
		}

		if strings.Contains(target, YerdBaseDir+"/") {
			return false, "", nil
		}

		return true, fmt.Sprintf("symlink to %s", target), nil
	}

	return true, "system binary", nil
}

// DetectSystemPHPInfo retrieves version information from system PHP installation.
// Returns PHP version string or error if version detection fails.
func DetectSystemPHPInfo() (string, error) {
	cmd := filepath.Join(SystemBinDir, "php")
	output, err := ExecuteCommand(cmd, "-v")
	if err != nil {
		return "", fmt.Errorf("failed to get PHP version: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "Unknown PHP version", nil
}

// GetPHPBinaryPath searches for PHP binary in multiple locations and validates version.
// version: PHP version string to locate. Returns binary path or error if not found.
func GetPHPBinaryPath(version string) (string, error) {
	possiblePaths := []string{
		YerdBinDir + "/php" + version,
		YerdPHPDir + "/php" + version + "/bin/php",
		SystemBinDir + "/php" + version,
		SystemBinDir + "/php",
		"/usr/bin/php" + version,
		"/usr/bin/php-" + version,
	}

	for _, path := range possiblePaths {
		if FileExists(path) {
			output, err := ExecuteCommand(path, "-v")
			if err == nil && strings.Contains(output, "PHP "+version) {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("PHP %s binary not found", version)
}

// GetPHPIniPath locates the php.ini configuration file for a specific PHP version.
// version: PHP version string. Returns ini file path or error if not found.
func GetPHPIniPath(version string) (string, error) {
	binaryPath, err := GetPHPBinaryPath(version)
	if err != nil {
		return "", err
	}

	output, err := ExecuteCommand(binaryPath, "--ini")
	if err != nil {
		return "", fmt.Errorf("failed to get ini info: %v", err)
	}

	lines := strings.Split(output, "\n")
	var configPath string
	var loadedIni string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Configuration File (php.ini) Path:") {
			configPath = strings.TrimSpace(strings.TrimPrefix(line, "Configuration File (php.ini) Path:"))
		} else if strings.HasPrefix(line, "Loaded Configuration File:") {
			loadedIni = strings.TrimSpace(strings.TrimPrefix(line, "Loaded Configuration File:"))
		}
	}

	if loadedIni != "" && loadedIni != "(none)" {
		return loadedIni, nil
	}

	if configPath != "" {
		potentialIniPath := configPath + phpIniFile
		if FileExists(potentialIniPath) {
			return potentialIniPath, nil
		}
		return fmt.Sprintf("%s (no php.ini)", configPath), nil
	}

	possibleIniPaths := []string{
		YerdEtcDir + "/php" + version + phpIniFile,
		YerdPHPDir + "/php" + version + "/lib/php.ini",
		YerdPHPDir + "/php" + version + "/etc/php.ini",
		"/etc/php/" + version + "/cli/php.ini",
		"/usr/local/etc/php/" + version + phpIniFile,
		"/opt/php/" + version + "/etc/php.ini",
	}

	for _, path := range possibleIniPaths {
		if FileExists(path) {
			return path, nil
		}
	}

	return "(none)", nil
}

// CreatePHPIniForVersion downloads and customizes a php.ini configuration file for a PHP version.
// version: PHP version string. Returns error if ini file already exists or creation fails.
func CreatePHPIniForVersion(version string) error {
	configDir := YerdEtcDir + "/php" + version
	iniPath := configDir + phpIniFile

	if FileExists(iniPath) {
		return fmt.Errorf("php.ini already exists at: %s", iniPath)
	}

	if err := os.MkdirAll(configDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	logger, err := NewLogger(fmt.Sprintf("php-ini-%s", version))
	if err != nil {
		return downloadAndCustomizePHPIni(version, iniPath, nil)
	}
	defer logger.DeleteLogFile()

	return downloadAndCustomizePHPIni(version, iniPath, logger)
}

// downloadAndCustomizePHPIni downloads the base php.ini and customizes it for the specific version
func downloadAndCustomizePHPIni(version, iniPath string, logger *Logger) error {
	if err := FetchConfigFromGitHub("php", "php.ini", iniPath, logger); err != nil {
		return fmt.Errorf("failed to download php.ini from GitHub: %v", err)
	}

	return customizePHPIni(version, logger)
}

// customizePHPIni reads the downloaded php.ini and updates version-specific settings
func customizePHPIni(version string, logger *Logger) error {
	extensionDir := detectExtensionDirectory(version)
	if extensionDir == "" {
		extensionDir = fmt.Sprintf("/opt/yerd/php/php%s/lib/php/extensions", version)
	}

	if err := UpdatePHPIniSetting(version, "extension_dir", extensionDir, logger); err != nil {
		return fmt.Errorf("failed to update extension_dir: %v", err)
	}

	SafeLog(logger, "Customized php.ini for PHP %s", version)
	return nil
}

// UpdatePHPIniSetting updates a specific setting in a PHP version's php.ini file
// It automatically detects the value type and formats it correctly (quoted strings vs unquoted values)
func UpdatePHPIniSetting(version, settingName, newValue string, logger *Logger) error {
	configDir := YerdEtcDir + "/php" + version
	iniPath := configDir + phpIniFile

	if !FileExists(iniPath) {
		return fmt.Errorf("php.ini not found for PHP %s at: %s", version, iniPath)
	}

	content, err := os.ReadFile(iniPath)
	if err != nil {
		return fmt.Errorf("failed to read php.ini: %v", err)
	}

	customizedContent := string(content)

	settingRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^(;?\s*%s\s*=\s*)(.*)$`, regexp.QuoteMeta(settingName)))

	if match := settingRegex.FindStringSubmatch(customizedContent); match != nil {
		originalValue := strings.TrimSpace(match[2])
		prefix := match[1]

		var formattedValue string
		if isQuotedValue(originalValue) {
			formattedValue = fmt.Sprintf(`"%s"`, newValue)
			SafeLog(logger, "Updating %s from %s to quoted value: %s", settingName, originalValue, formattedValue)
		} else {
			formattedValue = newValue
			SafeLog(logger, "Updating %s from %s to unquoted value: %s", settingName, originalValue, formattedValue)
		}

		cleanPrefix := regexp.MustCompile(`^;\s*`).ReplaceAllString(prefix, "")
		replacement := cleanPrefix + formattedValue
		customizedContent = settingRegex.ReplaceAllString(customizedContent, replacement)

		SafeLog(logger, "Updated %s setting in php.ini", settingName)
	} else {
		newLine := fmt.Sprintf("\n; Added by YERD\n%s = \"%s\"\n", settingName, newValue)
		customizedContent += newLine
		SafeLog(logger, "Added new %s setting to php.ini", settingName)
	}

	if err := os.WriteFile(iniPath, []byte(customizedContent), FilePermissions); err != nil {
		return fmt.Errorf("failed to write updated php.ini: %v", err)
	}

	SafeLog(logger, "Successfully updated %s = %s in PHP %s ini file", settingName, newValue, version)
	return nil
}

// isQuotedValue determines if a value is enclosed in quotes
func isQuotedValue(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) >= 2 &&
		((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\''))
}

// detectExtensionDirectory determines the correct extension directory path for PHP version.
// version: PHP version string. Returns extension directory path or empty string if not found.
func detectExtensionDirectory(version string) string {
	if extDir := getExtensionDirFromPHP(version); extDir != "" {
		return extDir
	}

	return findFallbackExtensionDir(version)
}

// getExtensionDirFromPHP queries PHP binary for its extension directory configuration.
// version: PHP version string. Returns extension directory from PHP or empty string.
func getExtensionDirFromPHP(version string) string {
	binaryPath, err := GetPHPBinaryPath(version)
	if err != nil {
		return ""
	}

	output, err := ExecuteCommand(binaryPath, "-r", "echo ini_get('extension_dir');")
	if err != nil {
		return ""
	}

	extDir := strings.TrimSpace(output)
	if extDir != "" && FileExists(extDir) {
		return extDir
	}

	return ""
}

// findFallbackExtensionDir searches common extension directory locations as fallback.
// version: PHP version string. Returns found extension directory or empty string.
func findFallbackExtensionDir(version string) string {
	installPath := YerdPHPDir + "/php" + version
	possibleExtDirs := []string{
		installPath + "/lib/php/extensions",
		installPath + "/lib64/php/extensions",
		installPath + "/lib/extensions",
	}

	for _, baseDir := range possibleExtDirs {
		if extDir := searchExtensionSubDir(baseDir); extDir != "" {
			return extDir
		}
	}

	return ""
}

// searchExtensionSubDir finds the actual extension directory within a base directory.
// baseDir: Base directory to search. Returns extension subdirectory path or baseDir.
func searchExtensionSubDir(baseDir string) string {
	if !FileExists(baseDir) {
		return ""
	}

	output, err := ExecuteCommand("find", baseDir, "-maxdepth", "1", "-type", "d", "-name", "*")
	if err != nil {
		return baseDir
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != baseDir && FileExists(line) {
			return line
		}
	}

	return baseDir
}

// CreateFPMMainConfig creates the main PHP-FPM configuration file for a specific PHP version.
// Downloads base configuration from GitHub and customizes it for the version.
// version: PHP version string. Returns error if main config creation fails.
func CreateFPMMainConfig(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	mainConfigDir := getFPMConfigDir(version)
	mainConfigPath := getFPMMainConfigFile(version)

	if FileExists(mainConfigPath) {
		return fmt.Errorf("FPM main config already exists at: %s", mainConfigPath)
	}

	if err := os.MkdirAll(mainConfigDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create FPM config directory %s: %v", mainConfigDir, err)
	}

	logger, err := NewLogger(fmt.Sprintf("php-fpm-main-%s", version))
	if err != nil {
		return downloadAndCustomizeFPMMain(version, mainConfigPath, nil)
	}
	defer func() {
		if logger != nil {
			logger.DeleteLogFile()
		}
	}()

	return downloadAndCustomizeFPMMain(version, mainConfigPath, logger)
}

// downloadAndCustomizeFPMMain downloads the main FPM config and customizes it for the version
func downloadAndCustomizeFPMMain(version, mainConfigPath string, logger *Logger) error {
	if err := FetchConfigFromGitHub("php", "php-fpm.conf", mainConfigPath, logger); err != nil {
		return fmt.Errorf("failed to download FPM main config from GitHub: %v", err)
	}

	return customizeFPMMain(version, mainConfigPath, logger)
}

// customizeFPMMain reads the downloaded main config and updates version-specific settings
func customizeFPMMain(version, mainConfigPath string, logger *Logger) error {
	if version == "" || mainConfigPath == "" {
		return fmt.Errorf("version and mainConfigPath cannot be empty")
	}

	data := TemplateData{
		"version":  version,
		"pid_path": filepath.Join(FPMSockDir, fmt.Sprintf("php%s-fpm.pid", version)),
		"log_path": filepath.Join(FPMLogDir, fmt.Sprintf("php%s-fpm.log", version)),
		"pool_dir": getFPMPoolConfigDir(version),
	}

	SafeLog(logger, "Reading FPM main config template from: %s", mainConfigPath)

	content, err := os.ReadFile(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read FPM main config template from %s: %v", mainConfigPath, err)
	}

	if len(content) == 0 {
		return fmt.Errorf("FPM main config template is empty")
	}

	SafeLog(logger, "Applying template substitutions for PHP %s main config", version)

	customizedContent := Template(string(content), data)

	if customizedContent == string(content) {
		SafeLog(logger, "Warning: No template substitutions were made for main config")
	}

	if err := os.WriteFile(mainConfigPath, []byte(customizedContent), FilePermissions); err != nil {
		return fmt.Errorf("failed to write customized FPM main config to %s: %v", mainConfigPath, err)
	}

	SafeLog(logger, "Successfully customized FPM main config for PHP %s", version)
	return nil
}

// CreateFPMPoolConfig creates a PHP-FPM pool configuration file for a specific PHP version.
// Downloads base configuration from GitHub and customizes it for the version.
// version: PHP version string. Returns error if pool config creation fails.
func CreateFPMPoolConfig(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	poolConfigDir := getFPMPoolConfigDir(version)
	poolConfigPath := getFPMPoolConfigFile(version)

	if FileExists(poolConfigPath) {
		return fmt.Errorf("FPM pool config already exists at: %s", poolConfigPath)
	}

	if err := os.MkdirAll(poolConfigDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create FPM config directory %s: %v", poolConfigDir, err)
	}

	logger, err := NewLogger(fmt.Sprintf("php-fpm-pool-%s", version))
	if err != nil {
		return downloadAndCustomizeFPMPool(version, poolConfigPath, nil)
	}
	defer func() {
		if logger != nil {
			logger.DeleteLogFile()
		}
	}()

	return downloadAndCustomizeFPMPool(version, poolConfigPath, logger)
}

// downloadAndCustomizeFPMPool downloads the base FPM pool config and customizes it for the version
func downloadAndCustomizeFPMPool(version, poolConfigPath string, logger *Logger) error {
	if err := FetchConfigFromGitHub("php", "www.conf", poolConfigPath, logger); err != nil {
		return fmt.Errorf("failed to download FPM pool config from GitHub: %v", err)
	}

	return customizeFPMPool(version, poolConfigPath, logger)
}

// customizeFPMPool reads the downloaded pool config and updates version-specific settings
func customizeFPMPool(version, poolConfigPath string, logger *Logger) error {
	if version == "" || poolConfigPath == "" {
		return fmt.Errorf("version and poolConfigPath cannot be empty")
	}

	data := TemplateData{
		"version":   version,
		"sock_path": filepath.Join(FPMSockDir, fmt.Sprintf("php%s-fpm.sock", version)),
		"log_path":  filepath.Join(FPMLogDir, fmt.Sprintf("php%s-fpm.log", version)),
		"user":      FPMUser,
		"group":     FPMGroup,
	}

	SafeLog(logger, "Reading FPM pool template from: %s", poolConfigPath)

	content, err := os.ReadFile(poolConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read FPM pool config template from %s: %v", poolConfigPath, err)
	}

	if len(content) == 0 {
		return fmt.Errorf("FPM pool config template is empty")
	}

	SafeLog(logger, "Applying template substitutions for PHP %s", version)

	customizedContent := Template(string(content), data)

	if customizedContent == string(content) {
		SafeLog(logger, "Warning: No template substitutions were made")
	}

	if err := os.WriteFile(poolConfigPath, []byte(customizedContent), FilePermissions); err != nil {
		return fmt.Errorf("failed to write customized FPM pool config to %s: %v", poolConfigPath, err)
	}

	SafeLog(logger, "Successfully customized FPM pool config for PHP %s with %s user/group", version, FPMUser)
	return nil
}

// SetupFPMEnvironment creates FPM directories and configuration files for a PHP version.
// version: PHP version string, forceRecreate: Whether to recreate existing configs, logger: Logger instance.
// Returns error if FPM setup fails.
func SetupFPMEnvironment(version string, forceRecreate bool, logger *Logger) error {
	SafeLog(logger, "Setting up FPM environment for PHP %s", version)

	// Create FPM runtime and log directories
	if err := CreateFPMDirectories(logger); err != nil {
		return err
	}

	// Create or recreate main config
	if err := CreateFPMConfigIfNeeded(version, "main", GetFPMMainConfigFile(version), CreateFPMMainConfig, forceRecreate, logger); err != nil {
		return err
	}

	// Create or recreate pool config
	if err := CreateFPMConfigIfNeeded(version, "pool", GetFPMPoolConfigFile(version), CreateFPMPoolConfig, forceRecreate, logger); err != nil {
		return err
	}

	SafeLog(logger, "FPM environment setup completed for PHP %s", version)
	return nil
}

// CreateFPMDirectories creates the required FPM runtime and log directories.
// logger: Logger instance. Returns error if directory creation fails.
func CreateFPMDirectories(logger *Logger) error {
	SafeLog(logger, "Creating FPM runtime directories")

	if err := CreateDirectory(FPMSockDir); err != nil {
		SafeLog(logger, "Failed to create FPM runtime directory: %v", err)
		return fmt.Errorf("failed to create FPM runtime directory: %v", err)
	}

	if err := CreateDirectory(FPMLogDir); err != nil {
		SafeLog(logger, "Failed to create FPM log directory: %v", err)
		return fmt.Errorf("failed to create FPM log directory: %v", err)
	}

	SafeLog(logger, "FPM directories created successfully")
	return nil
}

// CreateFPMConfigIfNeeded creates or recreates an FPM configuration file based on conditions.
// version: PHP version, configType: Type for logging, configPath: Path to config file,
// createFunc: Function to create config, forceRecreate: Whether to force recreation,
// logger: Logger instance. Returns error if config creation fails.
func CreateFPMConfigIfNeeded(version, configType, configPath string, createFunc func(string) error, forceRecreate bool, logger *Logger) error {
	configExists := FileExists(configPath)

	// Skip creation if file exists and not forcing recreation
	if configExists && !forceRecreate {
		SafeLog(logger, "FPM %s configuration already exists, skipping creation", configType)
		return nil
	}

	// Remove existing file if forcing recreation
	if forceRecreate && configExists {
		SafeLog(logger, "Recreating FPM %s configuration (force enabled)", configType)
		if err := os.Remove(configPath); err != nil {
			SafeLog(logger, "Warning: Failed to remove existing FPM %s config: %v", configType, err)
			// Continue despite removal failure
		}
	}

	// Create the configuration
	if err := createFunc(version); err != nil {
		SafeLog(logger, "Failed to create FPM %s config: %v", configType, err)
		return fmt.Errorf("failed to create FPM %s config: %v", configType, err)
	}

	SafeLog(logger, "FPM %s configuration created successfully", configType)
	return nil
}

// StartPHPFPM starts the PHP-FPM process for a specific PHP version using systemd.
// version: PHP version string. Returns error if FPM startup fails.
func StartPHPFPM(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	// Create systemd service if it doesn't exist
	servicePath := getSystemdServicePath(version)

	if !FileExists(servicePath) {
		if err := CreateSystemdService(version); err != nil {
			return fmt.Errorf("failed to create systemd service: %v", err)
		}
	}

	// Check if service is already active
	if IsSystemdServiceActive(version) {
		return fmt.Errorf("PHP-FPM %s is already running via systemd", version)
	}

	// Start the systemd service
	if err := StartSystemdService(version); err != nil {
		return fmt.Errorf("failed to start systemd service: %v", err)
	}

	// Wait a moment and verify it started
	time.Sleep(1000 * time.Millisecond)
	if !IsSystemdServiceActive(version) {
		return fmt.Errorf("PHP-FPM %s failed to start via systemd", version)
	}

	return nil
}

// IsProcessRunning checks if a process with the given PID is running.
// pid: Process ID as string. Returns true if process is running.
func IsProcessRunning(pid string) bool {
	pid = strings.TrimSpace(pid)
	if pid == "" {
		return false
	}

	_, err := ExecuteCommand("kill", "-0", pid)
	return err == nil
}

// StopPHPFPM stops the PHP-FPM process for a specific PHP version.
// version: PHP version string. Returns error if FPM stop fails.
func StopPHPFPM(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	// Check if service is running
	if !IsSystemdServiceActive(version) {
		return fmt.Errorf("PHP-FPM %s is not running via systemd", version)
	}

	// Stop the systemd service
	if err := StopSystemdService(version); err != nil {
		return fmt.Errorf("failed to stop systemd service: %v", err)
	}

	// Wait a moment and verify it stopped
	time.Sleep(1000 * time.Millisecond)
	if IsSystemdServiceActive(version) {
		return fmt.Errorf("PHP-FPM %s failed to stop via systemd", version)
	}

	return nil
}

// CreateSystemdService creates a systemd service file for PHP-FPM to ensure automatic startup on boot.
// version: PHP version string. Returns error if service creation fails.
func CreateSystemdService(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	servicePath := getSystemdServicePath(version)

	// Check if service already exists
	if FileExists(servicePath) {
		return fmt.Errorf("systemd service already exists at: %s", servicePath)
	}

	logger, err := NewLogger(fmt.Sprintf("systemd-service-%s", version))
	if err != nil {
		return createSystemdServiceFromGitHub(version, servicePath, nil)
	}
	defer func() {
		if logger != nil {
			logger.DeleteLogFile()
		}
	}()

	return createSystemdServiceFromGitHub(version, servicePath, logger)
}

// CreateSystemdServiceWithForce creates or recreates a systemd service for PHP-FPM with force option.
// version: PHP version string, force: Whether to overwrite existing service. Returns error if creation fails.
func CreateSystemdServiceWithForce(version string, force bool) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	servicePath := getSystemdServicePath(version)

	// If force is enabled and service exists, remove it first
	if force && FileExists(servicePath) {
		serviceName := getSystemdServiceName(version)
		// Stop service if running
		ExecuteCommand("systemctl", "stop", serviceName)
		// Remove existing service file
		if err := os.Remove(servicePath); err != nil {
			return fmt.Errorf("failed to remove existing systemd service file %s: %v", servicePath, err)
		}
	} else if FileExists(servicePath) && !force {
		return fmt.Errorf("systemd service already exists at: %s", servicePath)
	}

	logger, err := NewLogger(fmt.Sprintf("systemd-service-%s", version))
	if err != nil {
		return createSystemdServiceFromGitHub(version, servicePath, nil)
	}
	defer func() {
		if logger != nil {
			logger.DeleteLogFile()
		}
	}()

	return createSystemdServiceFromGitHub(version, servicePath, logger)
}

// createSystemdServiceFromGitHub downloads the systemd template and customizes it for the version
func createSystemdServiceFromGitHub(version, servicePath string, logger *Logger) error {
	// Download systemd template to temporary location
	tempPath := servicePath + ".tmp"
	if err := FetchConfigFromGitHub("php", "systemd.conf", tempPath, logger); err != nil {
		return fmt.Errorf("failed to download systemd service template from GitHub: %v", err)
	}
	defer os.Remove(tempPath) // Clean up temp file

	if err := customizeSystemdService(version, tempPath, servicePath, logger); err != nil {
		return err
	}

	serviceName := getSystemdServiceName(version)

	// Reload systemd daemon
	if _, err := ExecuteCommand("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd daemon: %v", err)
	}

	// Enable service for auto-start on boot
	if _, err := ExecuteCommand("systemctl", "enable", serviceName); err != nil {
		return fmt.Errorf("failed to enable systemd service %s: %v", serviceName, err)
	}

	SafeLog(logger, "Successfully created systemd service for PHP %s", version)
	return nil
}

// customizeSystemdService reads the downloaded systemd template and updates version-specific settings
func customizeSystemdService(version, tempPath, servicePath string, logger *Logger) error {
	fpmBinaryPath := filepath.Join(YerdPHPDir, "php"+version, "sbin", "php-fpm")
	mainConfigPath := getFPMMainConfigFile(version)
	pidPath := filepath.Join(FPMSockDir, "php"+version+"-fpm.pid")

	data := TemplateData{
		"version":          version,
		"fpm_binary_path":  fpmBinaryPath,
		"main_config_path": mainConfigPath,
		"pid_path":         pidPath,
	}

	SafeLog(logger, "Reading systemd service template from: %s", tempPath)

	content, err := os.ReadFile(tempPath)
	if err != nil {
		return fmt.Errorf("failed to read systemd service template from %s: %v", tempPath, err)
	}

	if len(content) == 0 {
		return fmt.Errorf("systemd service template is empty")
	}

	SafeLog(logger, "Applying template substitutions for PHP %s systemd service", version)

	customizedContent := Template(string(content), data)

	if customizedContent == string(content) {
		SafeLog(logger, "Warning: No template substitutions were made for systemd service")
	}

	// Write service file
	if err := WriteToFile(servicePath, []byte(customizedContent), FilePermissions); err != nil {
		return fmt.Errorf("failed to create systemd service file %s: %v", servicePath, err)
	}

	SafeLog(logger, "Successfully customized systemd service for PHP %s", version)
	return nil
}

// StartSystemdService starts the systemd service for PHP-FPM.
// version: PHP version string. Returns error if service start fails.
func StartSystemdService(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	serviceName := getSystemdServiceName(version)

	if _, err := ExecuteCommand("systemctl", "start", serviceName); err != nil {
		return fmt.Errorf("failed to start systemd service %s: %v", serviceName, err)
	}

	return nil
}

// StopSystemdService stops the systemd service for PHP-FPM.
// version: PHP version string. Returns error if service stop fails.
func StopSystemdService(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	serviceName := getSystemdServiceName(version)

	if _, err := ExecuteCommand("systemctl", "stop", serviceName); err != nil {
		return fmt.Errorf("failed to stop systemd service %s: %v", serviceName, err)
	}

	return nil
}

// RemoveSystemdService disables and removes the systemd service for PHP-FPM.
// version: PHP version string. Returns error if service removal fails.
func RemoveSystemdService(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	serviceName := getSystemdServiceName(version)
	servicePath := getSystemdServicePath(version)

	// Stop service if running
	ExecuteCommand("systemctl", "stop", serviceName)

	// Disable service
	ExecuteCommand("systemctl", "disable", serviceName)

	// Remove service file
	if FileExists(servicePath) {
		if err := os.Remove(servicePath); err != nil {
			return fmt.Errorf("failed to remove systemd service file %s: %v", servicePath, err)
		}
	}

	// Reload systemd daemon
	if _, err := ExecuteCommand("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd daemon: %v", err)
	}

	return nil
}

// IsSystemdServiceActive checks if the systemd service for PHP-FPM is active.
// version: PHP version string. Returns true if service is active.
func IsSystemdServiceActive(version string) bool {
	if err := validatePHPVersion(version); err != nil {
		return false
	}

	serviceName := getSystemdServiceName(version)
	_, err := ExecuteCommand("systemctl", "is-active", serviceName)
	return err == nil
}

// RemoveFPMPoolConfig removes the PHP-FPM pool configuration file for a specific PHP version.
// version: PHP version string. Returns error if removal fails.
func RemoveFPMPoolConfig(version string) error {
	if err := validatePHPVersion(version); err != nil {
		return err
	}

	poolConfigPath := getFPMPoolConfigFile(version)
	poolConfigDir := getFPMPoolConfigDir(version)

	// Remove pool config file
	if FileExists(poolConfigPath) {
		if err := os.Remove(poolConfigPath); err != nil {
			return fmt.Errorf("failed to remove FPM pool config %s: %v", poolConfigPath, err)
		}
	}

	// Remove php-fpm.d directory if empty
	if FileExists(poolConfigDir) {
		if err := os.Remove(poolConfigDir); err != nil {
			// Ignore error if directory is not empty
		}
	}

	return nil
}

// NormalizePHPVersion removes 'php' prefix from version strings for consistency.
// version: Version string potentially with php prefix. Returns normalized version string.
func NormalizePHPVersion(version string) string {
	if len(version) > 3 {
		prefix := strings.ToLower(version[:3])
		if prefix == "php" {
			return version[3:]
		}
	}
	return version
}
