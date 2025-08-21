package config

import (
	"fmt"
	"time"
)

type PhpInfo struct {
	Version          string    `json:"version"`
	InstalledVersion string    `json:"installed_version"`
	InstallPath      string    `json:"install_path"`
	InstallDate      time.Time `json:"install_date"`
	IsCLI            bool      `json:"is_cli"`
	Extensions       []string  `json:"extensions"`
	AddExtensions    []string  `json:"add_extensions"`
	RemoveExtensions []string  `json:"remove_extensions"`
}

type PhpConfig map[string]PhpInfo

// GetInstalledPhpInfo returns the PhpInfo struct for an installed
// version of php, however, if the version is not installed, then
// the return of bool will be false
// version: the version of php to check, eg: 8.1
func GetInstalledPhpInfo(version string) (PhpInfo, bool) {
	var info PhpInfo
	err := GetStruct(fmt.Sprintf("php.[%s]", version), info)
	if err == nil {
		return info, true
	}

	return PhpInfo{}, false
}
