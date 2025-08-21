package php

import (
	"fmt"
	"path/filepath"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

func UninstallPhp(info *config.PhpInfo) error {
	serviceName := fmt.Sprintf("yerd-php%s-fpm", info.Version)
	systemdPath := filepath.Join(constants.SystemdDir, fmt.Sprintf("yerd-php%s-fpm.service", info.Version))

	utils.SystemdStopService(serviceName)
	if err := utils.RemoveFile(systemdPath); err != nil {
		utils.LogError(err, "uninstall")
		return err
	}

	utils.SystemdReload()

	installDir := fmt.Sprintf("%s/php%s", constants.YerdPHPDir, info.Version)
	etcDir := fmt.Sprintf("%s/php%s", constants.YerdEtcDir, info.Version)
	localSymlink := fmt.Sprintf("%s/php%s", constants.YerdBinDir, info.Version)
	globalSymlink := fmt.Sprintf("%s/php%s", constants.SystemBinDir, info.Version)
	pidFile := fmt.Sprintf("%s/run/php%s-fpm.pid", constants.YerdPHPDir, info.Version)

	if err := utils.RunAll(
		func() error { return utils.RemoveSymlink(globalSymlink) },
		func() error { return utils.RemoveSymlink(localSymlink) },
		func() error { return utils.RemoveFile(pidFile) },
		func() error { return utils.RemoveFolder(installDir) },
		func() error { return utils.RemoveFolder(etcDir) },
		func() error {
			if info.IsCLI {
				return utils.RemoveSymlink(constants.GlobalPhpPath)
			}
			return nil
		},
	); err != nil {
		utils.LogError(err, "uninstall")
		utils.LogInfo("uninstall", "Failed to one one of the required functions")
		return err
	}

	config.Delete(fmt.Sprintf("php.[%s]", info.Version))

	return nil

}
