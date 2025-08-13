package cmd

import (
	"github.com/LumoSolutions/yerd/cmd/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Manage web services",
	Long:  `Install, configure, and manage web services like nginx, apache, and other web servers.`,
}

func init() {
	webCmd.AddCommand(web.InstallCmd)
}
