package cmd

import (
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yerd",
	Short: "A powerful, developer-friendly tool for managing PHP versions and local development environments with ease",
	Long: `Features:
  • Install and manage multiple PHP versions simultaneously
  • Switch PHP CLI versions instantly with simple commands
  • Lightweight and fast - no unnecessary overhead
  • Developer friendly`,
	Version: version.GetVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintSplash()
		cmd.Help()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(phpCmd)
	rootCmd.AddCommand(statusCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
