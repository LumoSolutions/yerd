package cmd

import (
	"github.com/spf13/cobra"
)

var composerCmd = &cobra.Command{
	Use:   "composer",
	Short: "Manage Composer",
	Long:  `Install, uninstall and update a YERD managed version of composer on your system.`,
}
