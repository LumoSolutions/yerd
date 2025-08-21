package cmd

import (
	"fmt"
	"os"

	"github.com/lumosolutions/yerd/cmd/composer"
	"github.com/lumosolutions/yerd/cmd/php"
	"github.com/lumosolutions/yerd/cmd/web"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/version"
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

// Execute runs the root command and handles any errors using cobra's error handler.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	phpCmd.AddCommand(php.BuildListCmd())
	phpCmd.AddCommand(php.BuildStatusCmd())
	phpVersions := constants.GetAvailablePhpVersions()
	for _, version := range phpVersions {
		phpCmd.AddCommand(php.CreateVersionCommand(version))
	}

	rootCmd.AddCommand(phpCmd)

	composerCmd.AddCommand(composer.BuildInstallCommand())
	composerCmd.AddCommand(composer.BuildUninstallCommand())
	composerCmd.AddCommand(composer.BuildUpdateCommand())

	rootCmd.AddCommand(composerCmd)

	webCmd.AddCommand(web.BuildInstallCommand())
	webCmd.AddCommand(web.BuildUninstallCommand())
	webCmd.AddCommand(web.BuildListCommand())
	webCmd.AddCommand(web.BuildAddCommand())
	webCmd.AddCommand(web.BuildRemoveCommand())
	webCmd.AddCommand(web.BuildSetCommand())

	rootCmd.AddCommand(webCmd)
}
