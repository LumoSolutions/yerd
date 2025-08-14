package web

import (
	"path/filepath"

	"github.com/LumoSolutions/yerd/internal/utils"
)

// FetchConfigFromGitHub downloads a configuration file from GitHub repository
func FetchConfigFromGitHub(service, configName string, logger *utils.Logger) error {
	configDir := GetServiceConfigPath(service)
	configPath := filepath.Join(configDir, configName)

	return utils.FetchConfigFromGitHub(service, configName, configPath, logger)
}
