package versions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/pkg/php"
)

type PHPReleaseResponse struct {
	Version           string       `json:"version"`
	Date              string       `json:"date"`
	Tags              []string     `json:"tags"`
	Source            []SourceFile `json:"source"`
	Announcement      bool         `json:"announcement"`
	SupportedVersions []string     `json:"supported_versions"`
}

type SourceFile struct {
	Filename string `json:"filename"`
	Name     string `json:"name"`
	SHA256   string `json:"sha256"`
	Date     string `json:"date"`
}

type VersionCache struct {
	LastUpdated    time.Time         `json:"last_updated"`
	LatestVersions map[string]string `json:"latest_versions"`
	DownloadURLs   map[string]string `json:"download_urls"`
}

const (
	PHPReleasesURL     = "https://www.php.net/releases/index.php?json&version="
	CacheValidDuration = 1 * time.Hour
	HTTPTimeout        = 10 * time.Second
	CacheFileExtension = "/version_cache.json"
	TempFileExtension  = ".tmp"
)

// FetchLatestVersions retrieves latest PHP version information from php.net for all supported versions.
// Returns latest versions map, download URLs map, or error if any version fetch fails.
func FetchLatestVersions() (map[string]string, map[string]string, error) {
	supportedMajorMinor := php.GetAvailableVersions()
	latestVersions := make(map[string]string)
	downloadURLs := make(map[string]string)

	for _, majorMinor := range supportedMajorMinor {
		latest, downloadURL, err := fetchLatestForMajorMinor(majorMinor)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch latest version for %s: %v", majorMinor, err)
		}
		latestVersions[majorMinor] = latest
		downloadURLs[latest] = downloadURL
	}

	return latestVersions, downloadURLs, nil
}

// fetchLatestForMajorMinor gets the latest version and download URL for a specific PHP major.minor version.
// majorMinor: Version string like "8.3". Returns version string, download URL, or error if fetch fails.
func fetchLatestForMajorMinor(majorMinor string) (string, string, error) {
	url := PHPReleasesURL + majorMinor

	client := &http.Client{Timeout: HTTPTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	var release PHPReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", fmt.Errorf("JSON decode failed: %v", err)
	}

	downloadURL := ""
	for _, source := range release.Source {
		if strings.HasSuffix(source.Filename, ".tar.gz") {
			downloadURL = fmt.Sprintf("https://www.php.net/distributions/%s", source.Filename)
			break
		}
	}

	if downloadURL == "" {
		return "", "", fmt.Errorf("no tar.gz download found for %s", majorMinor)
	}

	return release.Version, downloadURL, nil
}

// compareVersions compares two semantic version strings numerically.
// v1, v2: Version strings to compare. Returns 1 if v1>v2, -1 if v1<v2, 0 if equal.
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int

		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}

		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}

	return 0
}

// GetCachedVersions retrieves version information from local cache if valid.
// Returns cached version data and validity flag, or nil and false if cache invalid/missing.
func GetCachedVersions() (*VersionCache, bool) {
	configDir, err := utils.GetUserConfigDir()
	if err != nil {
		return nil, false
	}

	cacheFile := configDir + CacheFileExtension
	if !utils.FileExists(cacheFile) {
		return nil, false
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, false
	}

	var cache VersionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, false
	}

	if time.Since(cache.LastUpdated) > CacheValidDuration {
		return nil, false
	}

	return &cache, true
}

// SaveVersionCache writes version and download URL data to local cache file with proper ownership.
// latestVersions: Version mapping, downloadURLs: Download URL mapping. Returns error if save fails.
func SaveVersionCache(latestVersions, downloadURLs map[string]string) error {
	configDir, err := utils.GetUserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %v", err)
	}

	cache := VersionCache{
		LastUpdated:    time.Now(),
		LatestVersions: latestVersions,
		DownloadURLs:   downloadURLs,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %v", err)
	}

	cacheFile := configDir + CacheFileExtension

	tempFile := cacheFile + TempFileExtension
	if err := os.WriteFile(tempFile, data, utils.FilePermissions); err != nil {
		return fmt.Errorf("failed to write cache file: %v", err)
	}

	if err := os.Rename(tempFile, cacheFile); err != nil {
		return fmt.Errorf("failed to move cache file: %v", err)
	}

	userCtx, err := utils.GetRealUser()
	if err == nil && utils.CheckInstallPermissions() == nil {
		utils.ExecuteCommand("chown", fmt.Sprintf("%d:%d", userCtx.UID, userCtx.GID), cacheFile)
	}

	return nil
}

// GetLatestVersions retrieves PHP versions from cache or fetches fresh data if cache expired.
// Returns latest versions map, download URLs map, or error if fetch fails.
func GetLatestVersions() (map[string]string, map[string]string, error) {

	if cache, valid := GetCachedVersions(); valid {
		return cache.LatestVersions, cache.DownloadURLs, nil
	}

	latestVersions, downloadURLs, err := FetchLatestVersions()
	if err != nil {
		return nil, nil, err
	}

	if err := SaveVersionCache(latestVersions, downloadURLs); err != nil {

		fmt.Printf("Warning: failed to save version cache: %v\n", err)
	}

	return latestVersions, downloadURLs, nil
}

// GetLatestVersionsFresh bypasses cache and fetches fresh PHP version data from php.net.
// Returns latest versions map, download URLs map, or error if fetch fails.
func GetLatestVersionsFresh() (map[string]string, map[string]string, error) {

	latestVersions, downloadURLs, err := FetchLatestVersions()
	if err != nil {
		return nil, nil, err
	}

	if err := SaveVersionCache(latestVersions, downloadURLs); err != nil {

		fmt.Printf("Warning: failed to save version cache: %v\n", err)
	}

	return latestVersions, downloadURLs, nil
}

// CheckForUpdates compares installed PHP versions against latest available versions using cache.
// installedVersions: Map of installed version info. Returns update availability map or error.
func CheckForUpdates(installedVersions map[string]string) (map[string]bool, error) {
	latestVersions, _, err := GetLatestVersions()
	if err != nil {
		return nil, err
	}

	updates := make(map[string]bool)

	for majorMinor, installedVersion := range installedVersions {
		if latestVersion, exists := latestVersions[majorMinor]; exists {

			installedFull := extractVersionFromString(installedVersion)
			updates[majorMinor] = compareVersions(latestVersion, installedFull) > 0
		}
	}

	return updates, nil
}

// CheckForUpdatesFresh compares installed versions against fresh data and returns available updates.
// installedVersions: Map of installed version info. Returns update flags, available updates map, or error.
func CheckForUpdatesFresh(installedVersions map[string]string) (map[string]bool, map[string]string, error) {
	latestVersions, _, err := GetLatestVersionsFresh()
	if err != nil {
		return nil, nil, err
	}

	updates := make(map[string]bool)
	availableUpdates := make(map[string]string)

	for majorMinor, installedVersion := range installedVersions {
		if latestVersion, exists := latestVersions[majorMinor]; exists {

			installedFull := extractVersionFromString(installedVersion)
			hasUpdate := compareVersions(latestVersion, installedFull) > 0
			updates[majorMinor] = hasUpdate
			if hasUpdate {
				availableUpdates[majorMinor] = latestVersion
			}
		}
	}

	return updates, availableUpdates, nil
}

// extractVersionFromString extracts semantic version number from version string using regex.
// versionStr: Input version string. Returns extracted version or original string if no match.
func extractVersionFromString(versionStr string) string {

	versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+`)
	matches := versionRegex.FindString(versionStr)
	if matches != "" {
		return matches
	}
	return versionStr
}

// ExtractVersionFromString is a public wrapper for extractVersionFromString.
// versionStr: Input version string. Returns extracted semantic version number.
func ExtractVersionFromString(versionStr string) string {
	return extractVersionFromString(versionStr)
}
