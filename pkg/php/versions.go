package php

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
)

var availableVersions = []string{"8.1", "8.2", "8.3", "8.4"}

type VersionInfo struct {
	Version        string
	DownloadURL    string
	SourcePackage  string
	ConfigureFlags []string
}

func GetAvailableVersions() []string {
	return availableVersions
}

func IsValidVersion(version string) bool {
	for _, v := range availableVersions {
		if v == version {
			return true
		}
	}
	return false
}

func GetVersionInfo(version string) (VersionInfo, bool) {
	return VersionInfo{}, false
}

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
	
	// Get extension-specific configure flags
	extensionFlags := getExtensionConfigureFlags(extensions)
	
	return append(baseFlags, extensionFlags...)
}

func getExtensionConfigureFlags(extensions []string) []string {
	var flags []string
	
	// Extension to configure flag mapping
	extensionFlags := map[string]string{
		"mbstring":    "--enable-mbstring",
		"bcmath":      "--enable-bcmath",
		"opcache":     "--enable-opcache",
		"curl":        "--with-curl",
		"openssl":     "--with-openssl",
		"zip":         "--with-zip",
		"sockets":     "--enable-sockets",
		"mysqli":      "--with-mysqli",
		"pdo-mysql":   "--with-pdo-mysql",
		"gd":          "--enable-gd",
		"jpeg":        "--with-jpeg",
		"freetype":    "--with-freetype",
		"xml":         "--enable-xml",
		"json":        "--enable-json",
		"session":     "--enable-session",
		"hash":        "--enable-hash",
		"filter":      "--enable-filter",
		"pcre":        "--with-pcre-jit",
		"zlib":        "--with-zlib",
		"bz2":         "--with-bz2",
		"iconv":       "--with-iconv",
		"intl":        "--enable-intl",
		"pgsql":       "--with-pgsql",
		"pdo-pgsql":   "--with-pdo-pgsql",
		"sqlite3":     "--with-sqlite3",
		"pdo-sqlite":  "--with-pdo-sqlite",
		"fileinfo":    "--enable-fileinfo",
		"exif":        "--enable-exif",
		"gettext":     "--with-gettext",
		"gmp":         "--with-gmp",
		"ldap":        "--with-ldap",
		"soap":        "--enable-soap",
		"ftp":         "--enable-ftp",
	}
	
	for _, ext := range extensions {
		if flag, exists := extensionFlags[ext]; exists {
			flags = append(flags, flag)
		}
	}
	
	return flags
}

// Backward compatibility for existing code
func getConfigureFlagsForVersion(majorMinor string) []string {
	// Default extensions for backward compatibility
	defaultExtensions := []string{
		"mbstring", "bcmath", "opcache", "curl", "openssl", 
		"zip", "sockets", "mysqli", "pdo-mysql", "gd", "jpeg", "freetype",
	}
	return GetConfigureFlagsForVersion(majorMinor, defaultExtensions)
}

func GetLatestVersionInfo(majorMinor string, latestVersions, downloadURLs map[string]string) (VersionInfo, bool) {
	if latestVersion, exists := latestVersions[majorMinor]; exists {
		if downloadURL, hasURL := downloadURLs[latestVersion]; hasURL {
			return GetVersionInfoWithDownloadURL(majorMinor, downloadURL, latestVersion)
		}
	}

	return VersionInfo{}, false
}

func ExtractMajorMinor(fullVersion string) string {
	versionRegex := regexp.MustCompile(`^(\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(fullVersion)
	if len(matches) > 1 {
		return matches[1]
	}
	return fullVersion
}

func GetInstallPath(version string) string {
	return utils.YerdPHPDir + "/php" + version + "/"
}

func GetBinaryPath(version string) string {
	return utils.YerdBinDir + "/php" + version
}

func GetConfigPath(version string) string {
	return utils.YerdEtcDir + "/php" + version + "/"
}

func FormatVersion(version string) string {
	if len(version) > 3 {
		prefix := strings.ToLower(version[:3])
		if prefix == "php" {
			return version[3:]
		}
	}
	return version
}
