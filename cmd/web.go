package cmd

import (
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Manage Web Components",
	Long:  `Install, uninstall web components and manage local development sites.`,
}
