package php

import (
	"fmt"

	"github.com/fatih/color"
	phpinstaller "github.com/lumosolutions/yerd/internal/installers/php"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func buildCliCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cli",
		Short: fmt.Sprintf("Set default PHP CLI version to PHP %s", version),
		Long: `Set the default PHP version for command line usage.
		
Examples:
  yerd php 8.4 cli`,
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			if !utils.CheckAndPromptForSudo() {
				return
			}

			force, _ := cmd.Flags().GetBool("force")
			green := color.New(color.FgGreen)
			blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			data, installed := phpinstaller.IsInstalled(version)
			if !installed {
				red.Println("❌ Error: No action taken")
				blue.Printf("- PHP %s is not installed, please use\n", version)
				blue.Printf("- 'sudo yerd php %s install'\n\n", version)
				return
			}

			if data.IsCLI && !force {
				red.Println("❌ Error: No action taken")
				blue.Printf("- PHP %s is already the default CLI version of PHP\n", version)
				blue.Println("- If you wish to reapply this version forceably, you can use:")
				blue.Printf("- 'sudo yerd php %s cli -f'\n\n", version)
				return
			}

			fmt.Printf("Setting PHP %s as the default CLI version\n", version)

			if err := phpinstaller.SetCliVersion(data); err != nil {
				red.Println("❌ Error: No action taken")
				blue.Printf("- Unable to set PHP %s as the default CLI version", version)
				blue.Printf("- %v", err)
				return
			}

			green.Println("✓ Default PHP CLI version has been updated")
			blue.Printf("- PHP CLI version %s\n", version)

			if force {
				blue.Print("- update was forced using the -f/--force flag")
			}

		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force the regeneration of the symlinks for the php CLI")
	return cmd
}
