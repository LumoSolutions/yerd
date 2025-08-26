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

func buildUpdateCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: fmt.Sprintf("Update PHP %s", version),
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			if !utils.CheckAndPromptForSudo() {
				return
			}

			blue := color.New(color.FgBlue)
			yellow := color.New(color.FgYellow)
			red := color.New(color.FgRed)
			green := color.New(color.FgGreen)

			data, installed := config.GetInstalledPhpInfo(version)
			if !installed {
				yellow.Printf("PHP %s is not installed, please use the following command:\n", version)
				blue.Printf("- 'sudo yerd php %s install'\n\n", version)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			configFlag, _ := cmd.Flags().GetBool("config")
			versions, _, err := phpinstaller.GetLatestVersionsFresh()
			if err != nil {
				phpinstaller.PrintVersionFetchError(version)
				red.Printf("\n❌ Operation cancelled\n")
				return
			}

			if versions[version] == data.InstalledVersion {
				yellow.Printf("PHP %s is already running the latest version\n", version)
				blue.Printf("- Running version: %s\n", data.InstalledVersion)
				blue.Printf("- Latest version: %s\n", versions[version])
				blue.Printf("- To rebuild the current version, please use:\n")
				blue.Printf("- 'sudo yerd php %s rebuild\n\n", version)

				red.Printf("❌ Operation cancelled\n")
				return
			}

			installer, err := phpinstaller.NewPhpInstaller(version, true, configFlag)
			if err != nil {
				red.Printf("Failed to upgrade PHP %s: %v\n", version, err)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			if err := installer.Install(); err != nil {
				red.Printf("Failed to upgrade PHP %s: %v\n", version, err)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			green.Println("✓ PHP has been upgraded to the latest version successfully...")
			fmt.Println("Thanks for using YERD")
		},
	}

	cmd.Flags().BoolP("config", "c", false, "Recreate associated configuration, if it already exists")

	return cmd
}
