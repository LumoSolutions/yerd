package php

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fatih/color"
	phpinstaller "github.com/lumosolutions/yerd/internal/installers/php"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func buildUninstallCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: fmt.Sprintf("Uninstalls PHP %s", version),
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			if !utils.CheckAndPromptForSudo() {
				return
			}

			green := color.New(color.FgGreen)
			yellow := color.New(color.FgYellow)
			blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			agree, _ := cmd.Flags().GetBool("yes")

			data, installed := phpinstaller.IsInstalled(version)
			if !installed {
				red.Println("❌ Error: No action taken")
				blue.Printf("- PHP %s is not installed\n\n", version)
				return
			}

			if data.IsCLI && !agree {
				yellow.Printf("⚠️  Warning: PHP %s is currently set as CLI version\n", version)
				fmt.Printf("This will remove the PHP CLI and the 'php' command will no longer work.\n")
				fmt.Printf("Continue? (y/N): ")

				var response string
				fmt.Scanln(&response)

				if !isYes(response) {
					red.Printf("\n❌ Operation cancelled\n")
					return
				}

				fmt.Println()
			}

			if !agree {
				yellow.Printf("⚠️  Are you sure you want to uninstall PHP %s?\n", version)
				fmt.Printf("Confirm Action (y/N): ")
			}

			var response string
			fmt.Scanln(&response)

			if !isYes(response) {
				red.Printf("\n❌ Operation cancelled\n")
				return
			}

			fmt.Printf("Removing PHP %s\n", version)

			if err := phpinstaller.UninstallPhp(data); err != nil {
				red.Println("❌ Error: No action taken")
				blue.Printf("- Unable to uninstall PHP %s\n", version)
				blue.Printf("- %v\n\n", err)
				return
			}

			green.Printf("✓ PHP %s has been uninstalled\n\n", version)

		},
	}

	cmd.Flags().BoolP("agree", "y", false, "Automatically agree to any prompts")

	return cmd
}

func isYes(value string) bool {
	valid := []string{"y", "yes"}

	return slices.Contains(
		valid,
		strings.ToLower(strings.TrimSpace(value)),
	)
}
