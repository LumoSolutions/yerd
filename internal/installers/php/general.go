package php

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

// getPhpVersionInfo returns the version information for a PHP version
// version: PHP version to be fetched
// cached: Should we use the cache, or download fresh results from php.net
func getPhpVersionInfo(version string, cached bool) (*PhpVersionInfo, error) {
	var versions map[string]string
	var urls map[string]string
	var err error

	if cached {
		versions, urls, err = GetLatestVersions()
	} else {
		versions, urls, err = GetLatestVersionsFresh()
	}

	if err != nil {
		return nil, err
	}

	info, installed := config.GetInstalledPhpInfo(version)
	var extensions []string
	if installed {
		extensions = info.Extensions
	} else {
		extensions = constants.GetDefaultExtensions()
	}

	tempDir := os.TempDir()

	return &PhpVersionInfo{
		MajorMinor:     version,
		Version:        versions[version],
		DownloadURL:    urls[versions[version]],
		ConfigureFlags: getConfigureFlags(version, extensions),
		SourcePackage:  fmt.Sprintf("php-%s", version),
		ArchivePath:    filepath.Join(tempDir, fmt.Sprintf("php-%s.tar.gz", versions[version])),
		ExtractPath:    filepath.Join(tempDir, fmt.Sprintf("php-%s-extract", versions[version])),
		SourcePath:     filepath.Join(constants.YerdPHPDir, "src", fmt.Sprintf("php-%s", version)),
	}, nil
}

// GetFPMUser returns the username that should be used for PHP-FPM processes.
// Uses the real user context (handling sudo scenarios) or falls back to "nobody".
func GetFPMUser() string {
	userCtx, err := utils.GetRealUser()
	if err != nil {
		return "nobody"
	}
	return userCtx.Username
}

// GetFPMGroup returns the group name that should be used for PHP-FPM processes.
// Uses the real user's primary group (handling sudo scenarios) or falls back to "nobody".
func GetFPMGroup() string {
	userCtx, err := utils.GetRealUser()
	if err != nil {
		return "nobody"
	}

	// Get group information from GID
	group, err := user.LookupGroupId(strconv.Itoa(userCtx.GID))
	if err != nil {
		return "nobody"
	}

	return group.Name
}

func IsInstalled(version string) (*config.PhpInfo, bool) {
	var data *config.PhpInfo
	if err := config.GetStruct(fmt.Sprintf("php.[%s]", version), &data); err != nil {
		utils.LogError(err, "IsInstalled")
		return data, false
	}

	utils.LogInfo("IsInstalled", "Version %s", data.Version)
	return data, data.Version == version
}

func getBinaryPath(version string) string {
	return constants.YerdBinDir + "/php" + version
}
