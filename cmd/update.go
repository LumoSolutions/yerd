package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update YERD to the latest version",
	Long: `Check for and install the latest version of YERD from GitHub releases.

This command will:
- Check the current YERD version
- Fetch the latest release from GitHub
- Download and install the update if available
- Preserve existing configuration and PHP installations

Examples:
  yerd update           # Check for and install updates
  yerd update -y        # Auto-confirm update without prompting`,
	Args: cobra.NoArgs,
	Run:  runUpdate,
}

var autoConfirm bool

// runUpdate executes the YERD self-update process by checking for new releases and installing them.
func runUpdate(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	currentVersion := version.GetVersion()
	if currentVersion == "" {
		currentVersion = "unknown"
	}

	fmt.Printf("🔍 Current YERD version: %s\n", currentVersion)
	fmt.Printf("🌐 Checking for updates from GitHub...\n")

	latestRelease, err := fetchLatestRelease()
	if err != nil {
		fmt.Printf("❌ Failed to check for updates: %v\n", err)
		fmt.Printf("💡 Check your internet connection and try again\n")
		return
	}

	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
	if currentVersion == latestVersion {
		fmt.Printf("✅ YERD is already up to date (v%s)\n", currentVersion)
		return
	}

	if currentVersion != "unknown" && !isNewerVersion(latestVersion, currentVersion) {
		fmt.Printf("ℹ️  You have a newer or development version (v%s) than the latest release (v%s)\n", currentVersion, latestVersion)
		fmt.Printf("💡 No update needed\n")
		return
	}

	fmt.Printf("🆕 New version available: v%s\n", latestVersion)

	if !utils.CheckAndPromptForSudo() {
		return
	}

	if !confirmUpdate(latestVersion) {
		fmt.Printf("❌ Update cancelled\n")
		return
	}

	if err := performUpdate(latestRelease); err != nil {
		fmt.Printf("❌ Update failed: %v\n", err)
		fmt.Printf("💡 You can manually download from: https://github.com/LumoSolutions/yerd/releases\n")
		return
	}

	fmt.Printf("✅ YERD updated successfully to v%s\n", latestVersion)
	fmt.Printf("💡 Run 'yerd --version' to verify the update\n")
}

// fetchLatestRelease retrieves the latest YERD release information from GitHub API.
// Returns GitHubRelease struct with release details or error if request fails.
func fetchLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://api.github.com/repos/LumoSolutions/yerd/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %v", err)
	}

	return &release, nil
}

// isNewerVersion compares two semantic version strings to determine if latest is newer than current.
// latest: Version string to compare, current: Current version string. Returns true if latest > current.
func isNewerVersion(latest, current string) bool {
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	maxLen := len(latestParts)
	if len(currentParts) > maxLen {
		maxLen = len(currentParts)
	}

	for i := 0; i < maxLen; i++ {
		var latestNum, currentNum int

		if i < len(latestParts) {
			fmt.Sscanf(latestParts[i], "%d", &latestNum)
		}
		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentNum)
		}

		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	return false
}

// confirmUpdate prompts user for update confirmation or auto-confirms based on flags.
// version: Version string to update to. Returns true if user confirms or auto-confirm is enabled.
func confirmUpdate(version string) bool {
	if autoConfirm {
		fmt.Printf("🔄 Auto-updating to v%s...\n", version)
		return true
	}

	fmt.Printf("🔄 Update to v%s? (y/N): ", version)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}

// performUpdate downloads and installs the new YERD version from GitHub release.
// release: GitHubRelease containing download URLs and version info. Returns error if update fails.
func performUpdate(release *GitHubRelease) error {
	binaryName := getBinaryName()
	if binaryName == "" {
		return fmt.Errorf("no suitable binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("binary not found in release assets")
	}

	fmt.Printf("📦 Downloading %s...\n", binaryName)

	tempArchive, err := downloadBinary(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer os.Remove(tempArchive)

	fmt.Printf("📂 Extracting archive...\n")
	extractedBinary, err := extractBinaryFromArchive(tempArchive)
	if err != nil {
		return fmt.Errorf("extraction failed: %v", err)
	}
	defer os.Remove(extractedBinary)

	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %v", err)
	}

	fmt.Printf("🔄 Installing update...\n")

	backupPath := executablePath + ".backup"
	if err := os.Rename(executablePath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %v", err)
	}

	if err := copyFile(extractedBinary, executablePath); err != nil {
		os.Rename(backupPath, executablePath)
		return fmt.Errorf("failed to install new binary: %v", err)
	}

	if err := os.Chmod(executablePath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %v", err)
	}

	os.Remove(backupPath)

	return nil
}

// getBinaryName constructs the appropriate binary filename for the current platform and architecture.
// Returns platform-specific binary name or empty string if platform is unsupported.
func getBinaryName() string {
	version := strings.TrimPrefix(fetchLatestVersionTag(), "v")
	if version == "" {
		return ""
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goarch == "amd64" && goos == "linux" {
		return fmt.Sprintf("yerd_%s_linux_amd64.tar.gz", version)
	} else if goarch == "arm64" && goos == "linux" {
		return fmt.Sprintf("yerd_%s_linux_arm64.tar.gz", version)
	} else if goarch == "386" && goos == "linux" {
		return fmt.Sprintf("yerd_%s_linux_386.tar.gz", version)
	} else if goos == "darwin" {
		if goarch == "amd64" {
			return fmt.Sprintf("yerd_%s_darwin_amd64.tar.gz", version)
		} else if goarch == "arm64" {
			return fmt.Sprintf("yerd_%s_darwin_arm64.tar.gz", version)
		}
	}

	return ""
}

// fetchLatestVersionTag retrieves the latest version tag from GitHub releases.
// Returns version tag string or empty string if fetch fails.
func fetchLatestVersionTag() string {
	release, err := fetchLatestRelease()
	if err != nil {
		return ""
	}
	return release.TagName
}

// downloadBinary downloads a binary file from the given URL to a temporary file.
// url: Download URL. Returns path to temporary file or error if download fails.
func downloadBinary(url string) (string, error) {
	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "yerd-update-*.tar.gz")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

// extractBinaryFromArchive extracts the 'yerd' binary from a tar.gz archive.
// archivePath: Path to archive file. Returns path to extracted binary or error if extraction fails.
func extractBinaryFromArchive(archivePath string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if header.Name == "yerd" && header.Typeflag == tar.TypeReg {
			tempFile, err := os.CreateTemp("", "yerd-binary-*")
			if err != nil {
				return "", err
			}
			defer tempFile.Close()

			_, err = io.Copy(tempFile, tr)
			if err != nil {
				os.Remove(tempFile.Name())
				return "", err
			}

			return tempFile.Name(), nil
		}
	}

	return "", fmt.Errorf("yerd binary not found in archive")
}

// copyFile copies a file from src to dst, handling cross-device links properly.
// src: Source file path, dst: Destination file path. Returns error if copy fails.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
