package php

import (
	"fmt"

	phpinstaller "github.com/lumosolutions/yerd/internal/installers/php"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func buildRebuildCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: fmt.Sprintf("Rebuild PHP %s", version),
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			if !utils.CheckAndPromptForSudo() {
				return
			}

			data, installed := phpinstaller.IsInstalled(version)
			if !installed {
				fmt.Printf("PHP %s is not installed, please use\n", version)
				fmt.Printf("   'sudo yerd php %s install' instead\n\n", version)
				fmt.Println("Thanks for using YERD")
				return
			}

			nocache, _ := cmd.Flags().GetBool("nocache")
			config, _ := cmd.Flags().GetBool("config")

			if err := phpinstaller.RunRebuild(data, nocache, config); err != nil {
				fmt.Printf("Failed to rebuild php%s: %v\n", version, err)
				return
			}
		},
	}

	cmd.Flags().BoolP("nocache", "n", false, "Bypass cache to get the latest version from php.net")
	cmd.Flags().BoolP("config", "c", false, "Recreate associated configuration, if it already exists")

	return cmd
}
