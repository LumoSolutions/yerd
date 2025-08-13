package php

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/manager"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

var ComposerCmd = &cobra.Command{
	Use:   "composer",
	Short: "Install or update Composer",
	Long: `Install or update Composer PHP dependency manager.
	
Composer will be installed in the YERD directory structure and symlinked
to make it available globally.

Examples:
  yerd php composer    # Install or update Composer`,
	Args: cobra.NoArgs,
	Run:  runComposer,
}

// runComposer executes Composer installation/update with validation and error handling.
func runComposer(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo("Composer installation", "composer") {
		return
	}

	fmt.Println("Installing/Updating Composer...")

	err := manager.InstallComposer()
	if err != nil {
		fmt.Printf("\n❌ Composer installation failed: %v\n", err)
		fmt.Printf("💡 Run diagnostics: yerd doctor\n")
		return
	}

	fmt.Printf("✓ Composer installed/updated successfully\n")
	fmt.Printf("💡 You can now use 'composer' command globally\n")
}
