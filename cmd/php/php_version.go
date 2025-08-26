package php

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreateVersionCommand(version string) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   version,
		Short: fmt.Sprintf("Manage PHP %s", version),
		Long:  fmt.Sprintf("Commands for managing PHP version %s", version),
	}

	versionCmd.AddCommand(buildInstallCmd(version))
	versionCmd.AddCommand(buildRebuildCmd(version))
	versionCmd.AddCommand(buildExtensionsCmd(version))
	versionCmd.AddCommand(buildCliCmd(version))
	versionCmd.AddCommand(buildUninstallCmd(version))
	versionCmd.AddCommand(buildUpdateCmd(version))

	return versionCmd
}
