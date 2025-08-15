package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LumoSolutions/yerd/internal/version"
)

const (
	DefaultTimeout   = 30 * time.Second
	DefaultUserAgent = "YERD/1.0"
	ConfigRepoBase   = "https://raw.githubusercontent.com/LumoSolutions/yerd"
)

// DownloadOptions configures download behavior
type DownloadOptions struct {
	Timeout   time.Duration
	UserAgent string
	Logger    *Logger
}

// DefaultDownloadOptions returns default download configuration
func DefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		Timeout:   DefaultTimeout,
		UserAgent: DefaultUserAgent,
	}
}

// WithTimeout sets download timeout
func (opts *DownloadOptions) WithTimeout(timeout time.Duration) *DownloadOptions {
	opts.Timeout = timeout
	return opts
}

// WithUserAgent sets HTTP user agent
func (opts *DownloadOptions) WithUserAgent(userAgent string) *DownloadOptions {
	opts.UserAgent = userAgent
	return opts
}

// WithLogger sets logger for download operations
func (opts *DownloadOptions) WithLogger(logger *Logger) *DownloadOptions {
	opts.Logger = logger
	return opts
}

// DownloadFile downloads a file from URL to the specified path using http.Client with fallback to wget/curl
func DownloadFile(url, filePath string, opts *DownloadOptions) error {
	if opts == nil {
		opts = DefaultDownloadOptions()
	}

	SafeLog(opts.Logger, "Starting download from: %s", url)
	SafeLog(opts.Logger, "Destination: %s", filePath)

	dir := filepath.Dir(filePath)
	if dir != "" {
		if err := CreateDirectory(dir); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	if err := downloadWithHTTPClient(url, filePath, opts); err != nil {
		SafeLog(opts.Logger, "HTTP client failed: %v, trying command-line tools", err)

		// Fallback to command-line tools
		return downloadWithCommandLine(url, filePath, opts)
	}

	SafeLog(opts.Logger, "Download completed successfully")
	return nil
}

// downloadWithHTTPClient downloads using Go's http.Client
func downloadWithHTTPClient(url, filePath string, opts *DownloadOptions) error {
	client := &http.Client{
		Timeout: opts.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", opts.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(filePath)
		return fmt.Errorf("failed to write file: %v", err)
	}

	SafeLog(opts.Logger, "Downloaded using HTTP client")
	return nil
}

// downloadWithCommandLine downloads using wget or curl as fallback
func downloadWithCommandLine(url, filePath string, opts *DownloadOptions) error {
	if _, err := ExecuteCommand("which", "wget"); err == nil {
		SafeLog(opts.Logger, "Using wget for download")
		if _, err := ExecuteCommandWithLogging(opts.Logger, "wget", "-O", filePath, url); err != nil {
			os.Remove(filePath)
			return fmt.Errorf("wget download failed: %v", err)
		}
		return nil
	}

	// Try curl as fallback
	if _, err := ExecuteCommand("which", "curl"); err == nil {
		SafeLog(opts.Logger, "Using curl for download")
		if _, err := ExecuteCommandWithLogging(opts.Logger, "curl", "-L", "-o", filePath, url); err != nil {
			os.Remove(filePath)
			return fmt.Errorf("curl download failed: %v", err)
		}
		return nil
	}

	return fmt.Errorf("neither HTTP client, wget, nor curl are available for downloading")
}

// DownloadToTempDir downloads a file to a temporary directory and returns the path
func DownloadToTempDir(url, filename string, opts *DownloadOptions) (string, error) {
	tempDir := filepath.Join("/tmp", "yerd-download")
	if err := CreateDirectory(tempDir); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	tempPath := filepath.Join(tempDir, filename)
	if err := DownloadFile(url, tempPath, opts); err != nil {
		return "", err
	}

	return tempPath, nil
}

// DownloadAndExtractTarGz downloads and extracts a .tar.gz archive to the specified directory
func DownloadAndExtractTarGz(url, extractDir string, opts *DownloadOptions) (string, error) {
	filename := filepath.Base(url)
	if !strings.HasSuffix(filename, ".tar.gz") {
		filename += ".tar.gz"
	}

	tempPath, err := DownloadToTempDir(url, filename, opts)
	if err != nil {
		return "", fmt.Errorf("failed to download archive: %v", err)
	}
	defer os.Remove(tempPath)

	if err := CreateDirectory(extractDir); err != nil {
		return "", fmt.Errorf("failed to create extract directory: %v", err)
	}

	SafeLog(opts.Logger, "Extracting archive to: %s", extractDir)
	if _, err := ExecuteCommandWithLogging(opts.Logger, "tar", "xzf", tempPath, "-C", extractDir); err != nil {
		return "", fmt.Errorf("failed to extract archive: %v", err)
	}

	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return "", fmt.Errorf("failed to read extract directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			extractedPath := filepath.Join(extractDir, entry.Name())
			SafeLog(opts.Logger, "Archive extracted to: %s", extractedPath)
			return extractedPath, nil
		}
	}

	return extractDir, nil
}

// FetchTextContent downloads and returns the content of a text file as string
func FetchTextContent(url string, opts *DownloadOptions) (string, error) {
	if opts == nil {
		opts = DefaultDownloadOptions()
	}

	client := &http.Client{
		Timeout: opts.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", opts.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	SafeLog(opts.Logger, "Fetched text content from: %s", url)
	return string(content), nil
}

// FetchConfigFromGitHub downloads a configuration file from the GitHub repository to the specified path
func FetchConfigFromGitHub(configCategory, configName, destinationPath string, logger *Logger) error {
	return FetchConfigFromGitHubWithForce(configCategory, configName, destinationPath, false, logger)
}

// FetchConfigFromGitHubWithForce downloads a configuration file with option to force overwrite existing files
func FetchConfigFromGitHubWithForce(configCategory, configName, destinationPath string, force bool, logger *Logger) error {
	if !force && FileExists(destinationPath) {
		SafeLog(logger, "Config file already exists, skipping: %s", destinationPath)
		return nil
	}

	if force && FileExists(destinationPath) {
		SafeLog(logger, "Force overwriting existing config: %s", destinationPath)
	}

	configURL := fmt.Sprintf("%s/%s/.config/%s/%s", ConfigRepoBase, version.GetBranch(), configCategory, configName)
	SafeLog(logger, "Downloading config from: %s", configURL)

	opts := DefaultDownloadOptions().WithLogger(logger)
	if err := DownloadFile(configURL, destinationPath, opts); err != nil {
		return fmt.Errorf("failed to download config %s: %v", configName, err)
	}

	SafeLog(logger, "Downloaded config file: %s", destinationPath)
	return nil
}
