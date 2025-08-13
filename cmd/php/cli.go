package php

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/manager"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

var CliCmd = &cobra.Command{
	Use:   "cli [version]",
	Short: "Set default PHP CLI version",
	Long: `Set the default PHP version for command line usage.

Examples:
  yerd php cli 8.3
  yerd php cli 8.4
  yerd php cli php8.3`,
	Args: cobra.ExactArgs(1),
	Run:  runCli,
}

// runCli sets the default CLI PHP version with validation and system checks.
func runCli(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo("Setting CLI version", "cli", args[0]) {
		return
	}

	versionArg := args[0]
	phpVersion := utils.NormalizePHPVersion(versionArg)

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if _, exists := cfg.InstalledPHP[phpVersion]; !exists {
		fmt.Printf("Error: PHP %s is not installed\n", phpVersion)
		fmt.Printf("Install it first with: sudo yerd php add %s\n", phpVersion)
		return
	}

	fmt.Printf("Setting PHP %s as CLI version...\n", phpVersion)

	err = manager.SetCLIVersion(phpVersion)
	if err != nil {
		fmt.Printf("Error setting CLI version: %v\n", err)
		return
	}

	fmt.Printf("✓ PHP %s is now the default CLI version\n", phpVersion)
	fmt.Printf("Verify with: php -v\n")
}
