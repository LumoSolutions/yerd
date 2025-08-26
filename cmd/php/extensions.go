package php

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/config"
	phpinstaller "github.com/lumosolutions/yerd/internal/installers/php"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func buildExtensionsCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extensions <add|remove|list> <extension1> [extension2...]",
		Short: fmt.Sprintf("Manage PHP %s extensions", version),
		Long: fmt.Sprintf(`Manage PHP %s extensions.

Examples:
  yerd php %s extensions list                # List installed and available extensions
  yerd php %s extensions add gd memcached    # Add multiple extensions to PHP
  yerd php %s extensions remove gd           # Remove multiple extensions from PHP
  yerd php %s extensions add gd --rebuild    # Add extensions and automatically rebuild PHP`,
			version, version, version, version, version,
		),
		ValidArgs: []string{"add", "remove"},
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			if len(args) < 1 {
				red.Println("Error: requires at least 1 argument: <list|add|remove>")
				cmd.Usage()
				return
			}

			action := args[0]

			if len(args) < 2 && action != "list" {
				red.Printf("Error: requires at least 2 arguments: %s <extensions>\n", action)
				cmd.Usage()
				return
			}

			extensions := args[1:]

			rebuild, _ := cmd.Flags().GetBool("rebuild")
			nocache, _ := cmd.Flags().GetBool("nocache")
			configFlag, _ := cmd.Flags().GetBool("config")

			if rebuild {
				if !utils.CheckAndPromptForSudo() {
					return
				}
			}

			data, installed := config.GetInstalledPhpInfo(version)
			if !installed {
				red.Println("❌ Error: No action taken")
				blue.Printf("- PHP %s is not installed, please use\n", version)
				blue.Printf("- 'sudo yerd php %s install'\n\n", version)
				return
			}

			extManager := phpinstaller.NewExtensionManager(version, data, nocache, configFlag, rebuild)
			if err := extManager.RunAction(action, extensions); err != nil {
				return
			}
		},
	}

	cmd.Flags().BoolP("rebuild", "r", false, "Rebuild PHP after modifying extensions")
	cmd.Flags().BoolP("nocache", "n", false, "Bypass cache to get the latest version from php.net")
	cmd.Flags().BoolP("config", "c", false, "Recreate associated configuration, if it already exists")

	return cmd
}
