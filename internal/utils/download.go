package utils

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/lumosolutions/yerd/internal/constants"
)

const (
	DefaultTimeout   = 30 * time.Second
	DefaultUserAgent = "YERD/1.0"
	ConfigRepoBase   = "https://raw.githubusercontent.com/LumoSolutions/yerd"
)

type DownloadOptions struct {
	Timeout   time.Duration
	UserAgent string
}

// DefaultDownloadOptions returns default download configuration
func DefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		Timeout:   DefaultTimeout,
		UserAgent: DefaultUserAgent,
	}
}

// DownloadFile downloads a file to a given location
// url: The URL to download
// filePath: The location to place the file
// opts: Any download options required
func DownloadFile(url, filePath string, opts *DownloadOptions) error {
	if opts == nil {
		opts = DefaultDownloadOptions()
	}

	LogInfo(context, "Starting download from %s", url)
	LogInfo(context, "Destination: %s", filePath)

	dir := filepath.Dir(filePath)
	if dir != "" {
		if err := CreateDirectory(dir); err != nil {
			LogError(err, context)
			return fmt.Errorf("failed to create directory: %s", dir)
		}
	}

	if err := useHttpClient(url, filePath, opts); err != nil {
		LogWarning(context, "HTTP client failed %v, trying command-line tooling", err)
		if err := useCommandLine(url, filePath); err != nil {
			LogError(err, context)
			return fmt.Errorf("failed to download file")
		}
	}

	LogInfo(context, "File downloaded")
	return nil
}

func useHttpClient(url, filePath string, opts *DownloadOptions) error {
	client := &http.Client{
		Timeout: opts.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		LogError(err, context)
		return fmt.Errorf("failed to create request")
	}

	req.Header.Set("User-Agent", opts.UserAgent)
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if err := WriteToFile(filePath, body, constants.FilePermissions); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	LogInfo(context, "Downloaded using HTTP client")
	return nil
}

func useCommandLine(url, filePath string) error {
	if _, exists := CommandExists("wget"); exists {
		if _, success := ExecuteCommand("wget", "-O", filePath, url); !success {
			return fmt.Errorf("wget donwload failed")
		}

		LogInfo(context, "File downloaded using wget")
		return nil
	}

	if _, exists := CommandExists("curl"); exists {
		if _, success := ExecuteCommand("curl", "-L", "-o", filePath, url); !success {
			return fmt.Errorf("curl download failed")
		}

		LogInfo(context, "File downloaded using curl")
		return nil
	}

	return fmt.Errorf("wget & curl unavailable")
}
