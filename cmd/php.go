package cmd

import (
	"github.com/spf13/cobra"
	"github.com/LumoSolutions/yerd/cmd/php"
)

var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "Manage PHP versions",
	Long:  `Install, remove, update, and manage multiple PHP versions on your system.`,
}

func init() {
	phpCmd.AddCommand(php.AddCmd)
	phpCmd.AddCommand(php.RemoveCmd)
	phpCmd.AddCommand(php.CliCmd)
	phpCmd.AddCommand(php.ListCmd)
	phpCmd.AddCommand(php.DoctorCmd)
	phpCmd.AddCommand(php.UpdateCmd)
}