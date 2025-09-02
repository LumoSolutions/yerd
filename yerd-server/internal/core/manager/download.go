package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/lumosolutions/yerd/server/internal/constants"
	"github.com/lumosolutions/yerd/server/internal/core/services"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

type DownloadManager struct {
	Log     *services.Logger
	Timeout time.Duration
	client  *http.Client
}

func NewDownloadManager(log *services.Logger, timeout time.Duration) *DownloadManager {
	return &DownloadManager{
		Log:     log,
		Timeout: timeout,
	}
}

func (dm *DownloadManager) FetchJson(url string, target interface{}) error {
	body, err := dm.fetch(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode json: %v", err)
	}

	return nil
}

func (dm *DownloadManager) DownloadFile(url, path string) error {
	dm.Log.Info("download", "Downloading from %s", url)
	dm.Log.Info("download", "To path %s", path)

	directory := filepath.Dir(path)
	if !utils.IsDirectory(directory) {
		if err := utils.CreateDirectory(directory); err != nil {
			dm.Log.Error("download", err)
			return err
		}
	}

	client := dm.getClient()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		dm.Log.Error("download", err)
		return err
	}

	req.Header.Set("User-Agent", "YERD")
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")

	resp, err := client.Do(req)
	if err != nil {
		dm.Log.Error("download", err)
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	dm.Log.Info("download", "Status Code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		dm.Log.Error("download", err)
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if err := utils.WriteToFile(path, body, constants.FilePermissions); err != nil {
		dm.Log.Error("download", err)
		return fmt.Errorf("failed to create file: %v", err)
	}

	return nil
}

func (dm *DownloadManager) fetch(url string) ([]byte, error) {
	client := dm.getClient()

	dm.Log.Info("fetch", "Fetching: %s", url)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}

	defer resp.Body.Close()

	dm.Log.Info("fetch", "Response: %s", resp.StatusCode)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func (dm *DownloadManager) FetchFromGitHub(folder, file string) (string, error) {
	filePath := filepath.Join(".config", folder, file)

	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/%s/%s",
		constants.Repo,
		constants.Branch,
		filePath,
	)

	dm.Log.Info("github", "Attempting to download %s", url)

	content, err := dm.fetch(url)
	if err != nil {
		dm.Log.Error("download", err)
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}

	return string(content), nil
}

func (dm *DownloadManager) getClient() *http.Client {
	if dm.client == nil {
		dm.client = &http.Client{
			Timeout: dm.Timeout,
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		}
	}

	return dm.client
}
