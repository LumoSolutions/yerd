package php

import (
	"fmt"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

func SetCliVersion(info *config.PhpInfo) error {
	binaryPath := getBinaryPath(info.Version)
	globalPath := constants.GlobalPhpPath

	if err := utils.CreateSymlink(binaryPath, globalPath); err != nil {
		utils.LogError(err, "setcli")
		return err
	}

	info.IsCLI = true
	config.SetStruct(fmt.Sprintf("php.[%s]", info.Version), info)

	phpVersions := constants.GetAvailablePhpVersions()
	for _, version := range phpVersions {
		if version == info.Version {
			continue
		}

		if data, installed := config.GetInstalledPhpInfo(version); installed {
			if data.IsCLI {
				data.IsCLI = false
				config.SetStruct(fmt.Sprintf("php.[%s]", data.Version), data)
			}
		}
	}

	return nil
}
