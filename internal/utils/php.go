package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const phpIniFile = "/php.ini"

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

// CreateFPMPoolConfig creates a PHP-FPM pool configuration file for a specific PHP version.
// Downloads base configuration from GitHub and customizes it for the version.
// version: PHP version string. Returns error if pool config creation fails.
func CreateFPMPoolConfig(version string) error {
	if version == "" {
		return fmt.Errorf("PHP version cannot be empty")
	}

	configDir := filepath.Join(YerdEtcDir, "php"+version)
	poolConfigDir := filepath.Join(configDir, "php-fpm.d")
	poolConfigPath := filepath.Join(poolConfigDir, "www.conf")

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
		"pid_path":  filepath.Join(FPMPidDir, fmt.Sprintf("php%s-fpm.pid", version)),
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
