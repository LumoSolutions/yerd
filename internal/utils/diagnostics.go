package utils

import (
	"fmt"
	"strings"

	"github.com/LumoSolutions/yerd/internal/config"
)

type SystemRequirementsResult struct {
	BuildTools   map[string]bool
	Permissions  bool
	Commands     map[string]bool
	AllAvailable bool
}

type YERDConfigResult struct {
	Config         *config.Config
	ConfigError    error
	InstalledCount int
	CurrentCLI     string
	MissingDirs    []string
	ExistingDirs   []string
}

type SystemPHPResult struct {
	HasSystemPHP bool
	PHPType      string
	PHPInfo      string
	Error        error
}

type PHPVersionResult struct {
	Version     string
	IsValid     bool
	IsInstalled bool
	InstallDir  string
	BinaryPath  string
	VersionInfo string
	BinaryFound bool
}

type DirectoryStatus struct {
	Path        string
	Description string
	Exists      bool
}

// CheckSystemRequirements verifies availability of build tools, permissions, and commands for PHP compilation.
// Returns SystemRequirementsResult with tool availability status and overall system readiness.
func CheckSystemRequirements() *SystemRequirementsResult {
	result := &SystemRequirementsResult{
		BuildTools: make(map[string]bool),
		Commands:   make(map[string]bool),
	}

	buildTools := []string{"gcc", "make", "autoconf", "pkgconf", "bison", "re2c", "wget", "tar"}
	allBuildToolsAvailable := true

	for _, tool := range buildTools {
		if output, err := ExecuteCommand("which", tool); err == nil && output != "" {
			result.BuildTools[tool] = true
		} else {
			result.BuildTools[tool] = false
			allBuildToolsAvailable = false
		}
	}

	result.Permissions = CheckInstallPermissions() == nil

	commands := []string{"find", "which", "pacman"}
	for _, cmd := range commands {
		if output, err := ExecuteCommand("which", cmd); err == nil && output != "" {
			result.Commands[cmd] = true
		} else {
			result.Commands[cmd] = false
		}
	}

	result.AllAvailable = allBuildToolsAvailable
	return result
}

// CheckYERDConfiguration loads YERD config and validates directory structure.
// Returns YERDConfigResult with config status, installed versions, and directory information.
func CheckYERDConfiguration() *YERDConfigResult {
	result := &YERDConfigResult{}

	cfg, err := config.LoadConfig()
	result.Config = cfg
	result.ConfigError = err

	if err == nil {
		result.InstalledCount = len(cfg.InstalledPHP)
		result.CurrentCLI = cfg.CurrentCLI
	}

	dirs := []string{"/opt/yerd", "/opt/yerd/bin", "/opt/yerd/php", "/opt/yerd/etc"}
	for _, dir := range dirs {
		if FileExists(dir) {
			result.ExistingDirs = append(result.ExistingDirs, dir)
		} else {
			result.MissingDirs = append(result.MissingDirs, dir)
		}
	}

	return result
}

// CheckSystemPHPConflicts detects existing system PHP installations that may conflict with YERD.
// Returns SystemPHPResult with conflict status and PHP version information.
func CheckSystemPHPConflicts() *SystemPHPResult {
	result := &SystemPHPResult{}

	hasSystemPHP, phpType, err := CheckForSystemPHP()
	result.HasSystemPHP = hasSystemPHP
	result.PHPType = phpType
	result.Error = err

	if hasSystemPHP && err == nil {
		if info, err := DetectSystemPHPInfo(); err == nil {
			result.PHPInfo = info
		} else {
			result.PHPInfo = "Unknown PHP version"
		}
	}

	return result
}

// DiagnosePHPVersion performs comprehensive diagnostics on a specific PHP version.
// version: PHP version to diagnose, availableVersions: List of valid versions. Returns detailed diagnostic results.
func DiagnosePHPVersion(version string, availableVersions []string) *PHPVersionResult {
	result := &PHPVersionResult{
		Version: version,
	}

	result.IsValid = isValidPHPVersion(version, availableVersions)
	if !result.IsValid {
		return result
	}

	if cfg, err := config.LoadConfig(); err == nil {
		if _, exists := cfg.InstalledPHP[version]; exists {
			result.IsInstalled = true
		}
	}

	installDir := "/opt/yerd/php/php" + version
	result.InstallDir = installDir
	if FileExists(installDir) {

	}

	possiblePaths := []string{
		"/usr/bin/php" + version,
		"/usr/bin/php-" + version,
		"/usr/local/bin/php" + version,
	}

	for _, path := range possiblePaths {
		if FileExists(path) {
			result.BinaryPath = path
			result.BinaryFound = true

			if output, err := ExecuteCommand(path, "-v"); err == nil {
				lines := strings.Split(output, "\n")
				if len(lines) > 0 {
					result.VersionInfo = strings.TrimSpace(lines[0])
				}
			}
			break
		}
	}

	return result
}

// isValidPHPVersion checks if a version string exists in the available versions list.
// version: Version to validate, availableVersions: Valid versions list. Returns true if version is valid.
func isValidPHPVersion(version string, availableVersions []string) bool {
	for _, v := range availableVersions {
		if v == version {
			return true
		}
	}
	return false
}

// GetYERDDirectoryStatus checks existence of all required YERD directories.
// Returns slice of DirectoryStatus objects with path, description, and existence status.
func GetYERDDirectoryStatus() []DirectoryStatus {
	dirs := []DirectoryStatus{
		{"/opt/yerd", "YERD home", false},
		{"/opt/yerd/bin", "YERD binaries", false},
		{"/opt/yerd/php", "PHP installations", false},
		{"/opt/yerd/etc", "PHP configurations", false},
		{"/usr/local/bin", "System binaries", false},
	}

	for i := range dirs {
		dirs[i].Exists = FileExists(dirs[i].Path)
	}

	return dirs
}

// FindInstalledPHPBinaries searches common directories for PHP-related binaries.
// Returns map where keys are directory paths and values are lists of found PHP binaries.
func FindInstalledPHPBinaries() map[string][]string {
	result := make(map[string][]string)
	searchDirs := []string{"/usr/bin", "/usr/local/bin"}

	for _, dir := range searchDirs {
		var binaries []string
		output, err := ExecuteCommand("find", dir, "-name", "*php*", "-type", "f")
		if err == nil && output != "" {
			lines := strings.Split(strings.TrimSpace(output), "\n")
			for _, line := range lines {
				if line != "" {
					binaries = append(binaries, line)
				}
			}
		}
		result[dir] = binaries
	}

	return result
}

// PrintSystemRequirements displays formatted output of system requirements check results.
// result: SystemRequirementsResult containing tool and permission status.
func PrintSystemRequirements(result *SystemRequirementsResult) {

	buildTools := []string{"gcc", "make", "autoconf", "pkgconf", "bison", "re2c", "wget", "tar"}
	for _, tool := range buildTools {
		if available, exists := result.BuildTools[tool]; exists && available {
			fmt.Printf("├─ ✅ Build tool: %s (Available)\n", tool)
		} else {
			fmt.Printf("├─ ⚠️  Build tool: %s (Missing - will be auto-installed)\n", tool)
		}
	}

	if result.Permissions {
		fmt.Printf("├─ ✅ Permissions: Can write to system directories\n")
	} else {
		fmt.Printf("├─ ⚠️  Permissions: Requires sudo for installation\n")
	}

	commands := []string{"find", "which", "pacman"}
	for _, cmd := range commands {
		if available, exists := result.Commands[cmd]; exists && available {
			fmt.Printf("├─ ✅ Command available: %s\n", cmd)
		} else {
			fmt.Printf("├─ ❌ Command missing: %s\n", cmd)
		}
	}
}

// PrintYERDConfiguration displays formatted output of YERD configuration diagnostics.
// result: YERDConfigResult containing config status and directory information.
func PrintYERDConfiguration(result *YERDConfigResult) {
	if result.ConfigError != nil {
		fmt.Printf("├─ ❌ Config error: %v\n", result.ConfigError)
		return
	}

	fmt.Printf("├─ ✅ Config loaded: %d PHP versions installed\n", result.InstalledCount)

	if result.CurrentCLI != "" {
		fmt.Printf("├─ ✅ Current CLI: PHP %s\n", result.CurrentCLI)
	} else {
		fmt.Printf("├─ ℹ️  Current CLI: None set\n")
	}

	for _, dir := range result.ExistingDirs {
		fmt.Printf("├─ ✅ Directory exists: %s\n", dir)
	}
	for _, dir := range result.MissingDirs {
		fmt.Printf("├─ ⚠️  Directory missing: %s\n", dir)
	}
}

// PrintSystemPHPConflicts displays formatted output of system PHP conflict detection.
// result: SystemPHPResult containing conflict status and PHP information.
func PrintSystemPHPConflicts(result *SystemPHPResult) {
	if result.Error != nil {
		fmt.Printf("├─ ❌ Error checking system PHP: %v\n", result.Error)
		return
	}

	if result.HasSystemPHP {
		fmt.Printf("├─ ⚠️  System PHP detected: %s\n", result.PHPType)
		if result.PHPInfo != "" {
			fmt.Printf("├─ ℹ️  Version: %s\n", result.PHPInfo)
		}
		fmt.Printf("└─ 💡 Remove system PHP to use YERD CLI versions\n")
	} else {
		fmt.Printf("└─ ✅ No system PHP conflicts\n")
	}
}

// PrintPHPVersionDiagnostics displays detailed diagnostic information for a specific PHP version.
// result: PHPVersionResult with diagnostic data, availableVersions: Valid versions list.
func PrintPHPVersionDiagnostics(result *PHPVersionResult, availableVersions []string) {
	if !result.IsValid {
		fmt.Printf("├─ ❌ Invalid PHP version: %s\n", result.Version)
		fmt.Printf("└─ Valid versions: %s\n", strings.Join(availableVersions, ", "))
		return
	}

	if result.IsInstalled {
		fmt.Printf("├─ ✅ YERD status: Installed\n")
	} else {
		fmt.Printf("├─ ❌ YERD status: Not installed\n")
	}

	if FileExists(result.InstallDir) {
		fmt.Printf("├─ ✅ Source installation: PHP %s installed\n", result.Version)
	} else {
		fmt.Printf("├─ ❌ Source installation: PHP %s not installed\n", result.Version)
	}

	if result.BinaryFound {
		fmt.Printf("├─ ✅ Binary found: %s\n", result.BinaryPath)
		if result.VersionInfo != "" {
			fmt.Printf("└─ ℹ️  Version info: %s\n", result.VersionInfo)
		}
	} else {
		fmt.Printf("└─ ❌ Binary not found in common locations\n")
	}
}

// PrintInstalledPHPBinaries displays found PHP binaries organized by directory location.
// binaries: Map of directory paths to PHP binary file lists.
func PrintInstalledPHPBinaries(binaries map[string][]string) {
	for dir, files := range binaries {
		fmt.Printf("📁 %s:\n", dir)
		if len(files) > 0 {
			for _, file := range files {
				fmt.Printf("  - %s\n", file)
			}
		} else {
			fmt.Printf("  - No PHP binaries found\n")
		}
	}
}
