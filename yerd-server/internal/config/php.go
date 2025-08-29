package config

import (
	"time"
)

type PhpConfig map[string]PhpInfo

type PhpInfo struct {
	Version          string    `json:"version"`
	InstalledVersion string    `json:"installed_version"`
	InstallPath      string    `json:"install_path"`
	InstallDate      time.Time `json:"install_date"`
	IsCLI            bool      `json:"is_cli"`
	Global           bool      `json:"global"`
	Extensions       []string  `json:"extensions"`
	NeedsRebuild     bool      `json:"needs_rebuild"`
	FpmPidLocation   string    `json:"fpm_pid_location"`
	FpmSocket        string    `json:"fpm_socket"`
	PhpIniLocation   string    `json:"php_ini_location"`
	FpmConfig        string    `json:"fpm_config"`
}
