package php

import (
	"fmt"

	"github.com/fatih/color"
	phpinstaller "github.com/lumosolutions/yerd/internal/installers/php"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func buildInstallCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: fmt.Sprintf("Install PHP %s", version),
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			if !utils.CheckAndPromptForSudo() {
				return
			}

			blue := color.New(color.FgBlue)
			yellow := color.New(color.FgYellow)
			red := color.New(color.FgRed)
			green := color.New(color.FgGreen)

			if _, installed := phpinstaller.IsInstalled(version); installed {
				yellow.Printf("PHP %s is already installed, please use one of the following:\n", version)
				blue.Printf("- 'sudo yerd php %s rebuild' to build the current version\n", version)
				blue.Printf("- 'sudo yerd php %s upgrade' to update PHP %s to the latest version\n\n", version, version)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			nocache, _ := cmd.Flags().GetBool("nocache")

			installer, err := phpinstaller.NewPhpInstaller(version, nocache, true)
			if err != nil {
				red.Printf("Failed to install php%s: %v\n", version, err)
				return
			}

			if err := installer.Install(); err != nil {
				red.Printf("Failed to install php%s: %v\n", version, err)
				return
			}

			green.Println("✓ Installation complete...")
			fmt.Println("Thanks for using YERD")
		},
	}

	cmd.Flags().BoolP("nocache", "n", false, "Bypass cache to get the latest version from php.net")

	return cmd
}
