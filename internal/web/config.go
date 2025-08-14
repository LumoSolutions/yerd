package web

import (
	"fmt"
	"path/filepath"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
)

const (
	ConfigRepoBase = "https://raw.githubusercontent.com/LumoSolutions/yerd"
)

// FetchConfigFromGitHub downloads a configuration file from GitHub repository
func FetchConfigFromGitHub(service, configName string, logger *utils.Logger) error {
	configDir := GetServiceConfigPath(service)
	configPath := filepath.Join(configDir, configName)

	if utils.FileExists(configPath) {
		logger.WriteLog("Config file already exists, skipping: %s", configPath)
		return nil
	}

	configURL := fmt.Sprintf("%s/%s/.config/%s/%s", ConfigRepoBase, version.GetBranch(), service, configName)
	logger.WriteLog("Downloading config from: %s", configURL)

	opts := utils.DefaultDownloadOptions().WithLogger(logger)
	if err := utils.DownloadFile(configURL, configPath, opts); err != nil {
		return fmt.Errorf("failed to download config %s: %v", configName, err)
	}

	logger.WriteLog("Downloaded config file: %s", configPath)
	return nil
}
