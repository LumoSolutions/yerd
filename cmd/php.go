package cmd

import (
	"github.com/spf13/cobra"
)

var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "Manage PHP versions",
	Long:  `Install, remove, update, and manage multiple PHP versions on your system.`,
}
