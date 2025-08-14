package php

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/manager"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove [version]",
	Short: "Remove a specific PHP version",
	Long: `Remove an installed PHP version and clean up symlinks.

Examples:
  yerd php remove 8.3
  yerd php remove 8.4
  yerd php remove php8.3`,
	Args: cobra.ExactArgs(1),
	Run:  runRemove,
}

// runRemove executes PHP version removal with CLI version safety checks and user confirmation.
func runRemove(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo("PHP removal", "remove", args[0]) {
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
		return
	}

	if cfg.CurrentCLI == phpVersion {
		fmt.Printf("⚠️  Warning: PHP %s is currently set as CLI version\n", phpVersion)
		fmt.Printf("This will remove the CLI symlink and 'php' command will no longer work.\n")
		fmt.Printf("Continue? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
			fmt.Printf("❌ Operation cancelled\n")
			return
		}
		fmt.Println()
	}

	fmt.Printf("Removing PHP %s...\n", phpVersion)

	// Stop and remove systemd service first
	if utils.IsSystemdServiceActive(phpVersion) {
		fmt.Printf("Stopping PHP-FPM service...\n")
		if err := utils.StopPHPFPM(phpVersion); err != nil {
			utils.PrintWarning("Failed to stop PHP-FPM service: %v", err)
		}
	}

	fmt.Printf("Removing systemd service...\n")
	if err := utils.RemoveSystemdService(phpVersion); err != nil {
		utils.PrintWarning("Failed to remove systemd service: %v", err)
	}

	fmt.Printf("Removing FPM pool configuration...\n")
	if err := utils.RemoveFPMPoolConfig(phpVersion); err != nil {
		utils.PrintWarning("Failed to remove FPM pool config: %v", err)
	}

	err = manager.RemovePHP(phpVersion)
	if err != nil {
		fmt.Printf("Error removing PHP %s: %v\n", phpVersion, err)
		return
	}

	fmt.Printf("✓ PHP %s removed successfully\n", phpVersion)
}
