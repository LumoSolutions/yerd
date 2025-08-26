package cmd

import (
	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manage local development sites",
	Long:  `Add, remove and manage local development sites.`,
}
