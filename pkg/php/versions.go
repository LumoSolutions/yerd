package php

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/pkg/extensions"
)

var availableVersions = []string{"8.1", "8.2", "8.3", "8.4"}

type VersionInfo struct {
	Version        string
	DownloadURL    string
	SourcePackage  string
	ConfigureFlags []string
}

// GetAvailableVersions returns the list of PHP versions supported by YERD.
func GetAvailableVersions() []string {
	return availableVersions
}

// IsValidVersion checks if the provided version string is supported by YERD.
// version: PHP version string to validate. Returns true if version is supported.
func IsValidVersion(version string) bool {
	for _, v := range availableVersions {
		if v == version {
			return true
		}
	}
	return false
}

// GetVersionInfo retrieves version information for a PHP version (placeholder implementation).
// version: PHP version string. Returns VersionInfo struct and existence boolean.
func GetVersionInfo(version string) (VersionInfo, bool) {
	return VersionInfo{}, false
}

// GetVersionInfoWithDownloadURL creates VersionInfo with download details and configure flags.
// version: Major.minor version, downloadURL: Source download URL, fullVersion: Complete version string.
func GetVersionInfoWithDownloadURL(version, downloadURL, fullVersion string) (VersionInfo, bool) {
	if downloadURL != "" && fullVersion != "" {
		return VersionInfo{
			Version:        version,
			DownloadURL:    downloadURL,
			SourcePackage:  fmt.Sprintf("php-%s", fullVersion),
			ConfigureFlags: getConfigureFlagsForVersion(version),
		}, true
	}

	return VersionInfo{}, false
}

// GetConfigureFlagsForVersion generates PHP configure flags for compilation with specified extensions.
// majorMinor: PHP version, extensions: Extension list. Returns complete configure flag slice.
func GetConfigureFlagsForVersion(majorMinor string, extensions []string) []string {
	baseFlags := []string{
		fmt.Sprintf("--prefix=%s/php%s", utils.YerdPHPDir, majorMinor),
		fmt.Sprintf("--with-config-file-path=%s/php%s", utils.YerdEtcDir, majorMinor),
		fmt.Sprintf("--with-config-file-scan-dir=%s/php%s/conf.d", utils.YerdEtcDir, majorMinor),
		"--enable-fpm",
		"--with-fpm-user=http",
		"--with-fpm-group=http",
		"--enable-cli",
	}

	extensionFlags := getExtensionConfigureFlags(extensions)

	return append(baseFlags, extensionFlags...)
}

// getExtensionConfigureFlags converts extension names to their configure flags.
// extensions: Extension name list. Returns slice of configure flags for enabled extensions.
func getExtensionConfigureFlags(extList []string) []string {
	return extensions.GetConfigureFlags(extList)
}

// getConfigureFlagsForVersion generates configure flags with default extensions for a PHP version.
// majorMinor: PHP version string. Returns configure flags with standard extension set.
func getConfigureFlagsForVersion(majorMinor string) []string {
	defaultExtensions := []string{
		"mbstring", "bcmath", "opcache", "curl", "openssl",
		"zip", "sockets", "mysqli", "pdo-mysql", "gd", "jpeg", "freetype",
	}
	return GetConfigureFlagsForVersion(majorMinor, defaultExtensions)
}

// GetLatestVersionInfo creates VersionInfo from version and URL mappings.
// majorMinor: PHP version, latestVersions: Version mapping, downloadURLs: URL mapping. Returns VersionInfo or false.
func GetLatestVersionInfo(majorMinor string, latestVersions, downloadURLs map[string]string) (VersionInfo, bool) {
	if latestVersion, exists := latestVersions[majorMinor]; exists {
		if downloadURL, hasURL := downloadURLs[latestVersion]; hasURL {
			return GetVersionInfoWithDownloadURL(majorMinor, downloadURL, latestVersion)
		}
	}

	return VersionInfo{}, false
}

// ExtractMajorMinor extracts major.minor version from full version string using regex.
// fullVersion: Complete version string like "8.3.1". Returns major.minor format like "8.3".
func ExtractMajorMinor(fullVersion string) string {
	versionRegex := regexp.MustCompile(`^(\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(fullVersion)
	if len(matches) > 1 {
		return matches[1]
	}
	return fullVersion
}

// GetInstallPath returns the installation directory path for a PHP version.
// version: PHP version string. Returns full path to PHP installation directory.
func GetInstallPath(version string) string {
	return utils.YerdPHPDir + "/php" + version + "/"
}

// GetBinaryPath returns the symlink path for a PHP version's binary.
// version: PHP version string. Returns path to PHP binary symlink in YERD bin directory.
func GetBinaryPath(version string) string {
	return utils.YerdBinDir + "/php" + version
}

// GetConfigPath returns the configuration directory path for a PHP version.
// version: PHP version string. Returns path to PHP configuration directory.
func GetConfigPath(version string) string {
	return utils.YerdEtcDir + "/php" + version + "/"
}

// FormatVersion removes 'php' prefix from version strings for normalization.
// version: Version string potentially with php prefix. Returns normalized version string.
func FormatVersion(version string) string {
	if len(version) > 3 {
		prefix := strings.ToLower(version[:3])
		if prefix == "php" {
			return version[3:]
		}
	}
	return version
}
