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

func getConfigureFlagsForVersion(majorMinor string) []string {
	return []string{
		fmt.Sprintf("--prefix=%s/php%s", utils.YerdPHPDir, majorMinor),
		fmt.Sprintf("--with-config-file-path=%s/php%s", utils.YerdEtcDir, majorMinor),
		fmt.Sprintf("--with-config-file-scan-dir=%s/php%s/conf.d", utils.YerdEtcDir, majorMinor),
		"--enable-fpm",
		"--with-fpm-user=http",
		"--with-fpm-group=http",
		"--enable-cli",
		"--enable-mbstring",
		"--enable-bcmath",
		"--enable-opcache",
		"--with-curl",
		"--with-openssl",
		"--with-zip",
		"--enable-sockets",
		"--with-mysqli",
		"--with-pdo-mysql",
		"--enable-gd",
		"--with-jpeg",
		"--with-freetype",
	}
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
