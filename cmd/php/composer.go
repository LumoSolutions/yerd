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
	Short: "Install, update, or remove Composer",
	Long: `Install, update, or remove Composer PHP dependency manager.
	
Composer will be installed in the YERD directory structure and symlinked
to make it available globally.

Examples:
  yerd php composer       # Install or update Composer
  yerd php composer -r    # Remove Composer
  yerd php composer --remove  # Remove Composer`,
	Args: cobra.NoArgs,
	Run:  runComposer,
}

func init() {
	ComposerCmd.Flags().BoolP("remove", "r", false, "Remove Composer installation")
}

// runComposer executes Composer installation/update or removal with validation and error handling.
func runComposer(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	remove, _ := cmd.Flags().GetBool("remove")

	if remove {
		handleComposerRemoval()
	} else {
		handleComposerInstallation()
	}
}

// handleComposerInstallation manages the composer installation process.
func handleComposerInstallation() {
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

// handleComposerRemoval manages the composer removal process.
func handleComposerRemoval() {
	if !utils.CheckAndPromptForSudo("Composer removal", "composer", "--remove") {
		return
	}

	fmt.Println("Removing Composer...")

	err := manager.RemoveComposer()
	if err != nil {
		fmt.Printf("\n❌ Composer removal failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Composer removed successfully\n")
}
