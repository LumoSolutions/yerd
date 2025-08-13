package php

import (
	"fmt"
	"strings"

	"github.com/LumoSolutions/yerd/internal/builder"
	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/pkg/constants"
	"github.com/LumoSolutions/yerd/pkg/php"
	"github.com/spf13/cobra"
)

var RebuildCmd = &cobra.Command{
	Use:   "rebuild <version>",
	Short: "Rebuild PHP with current extensions",
	Long: `Force rebuild PHP with the currently configured extensions.
	
This is useful for:
  - Recovering from corrupted PHP installations
  - Applying system updates to dependencies
  - Troubleshooting build issues
  - Ensuring extensions are properly compiled`,
	Args: cobra.ExactArgs(1),
	RunE: runRebuild,
}

// runRebuild forces a complete rebuild of PHP with existing extensions configuration.
// Returns error if rebuild fails, nil if successful.
func runRebuild(cmd *cobra.Command, args []string) error {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo("PHP rebuild", "rebuild", args[0]) {
		return nil
	}

	phpVersion := php.FormatVersion(args[0])

	if !php.IsValidVersion(phpVersion) {
		utils.PrintError("Invalid PHP version: %s", phpVersion)
		return nil
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.PrintError("Failed to load config: %v", err)
		return nil
	}

	if _, exists := cfg.InstalledPHP[phpVersion]; !exists {
		utils.PrintError("PHP %s is not installed. Use 'yerd php add %s' first", phpVersion, phpVersion)
		return nil
	}

	currentExtensions, exists := cfg.GetPHPExtensions(phpVersion)
	if !exists {
		utils.PrintError("No extension information found for PHP %s", phpVersion)
		return nil
	}

	if len(currentExtensions) == 0 {
		utils.PrintWarning("PHP %s has no extensions configured. Adding default extensions.", phpVersion)
		currentExtensions = constants.DefaultPHPExtensions
		cfg.UpdatePHPExtensions(phpVersion, currentExtensions)
	}

	utils.PrintInfo("Rebuilding PHP %s with extensions: %s", phpVersion, strings.Join(currentExtensions, ", "))

	if err := forceRebuildPHP(cfg, phpVersion, currentExtensions); err != nil {
		return nil
	}

	return nil
}

// forceRebuildPHP performs the actual rebuild process with spinner animation.
// cfg: Configuration object, version: PHP version to rebuild, extensions: Extensions to include.
func forceRebuildPHP(cfg *config.Config, version string, extensions []string) error {
	utils.PrintWarning("Force rebuilding PHP (no configuration backup needed)...")

	spinner := utils.NewLoadingSpinner(fmt.Sprintf("Building PHP %s with extensions", version))
	spinner.Start()

	phpBuilder, err := builder.NewBuilder(version, extensions)
	if err != nil {
		spinner.Stop("✗ Failed to create builder")
		utils.PrintError("Failed to create builder: %v", err)
		return fmt.Errorf("builder creation failed")
	}
	err = phpBuilder.RebuildPHP()

	if err != nil {
		spinner.Stop("✗ Build failed")
		utils.PrintError("Failed to rebuild PHP %s: %v", version, err)
		phpBuilder.Cleanup()
		return fmt.Errorf("rebuild failed")
	}

	spinner.Stop("✓ Build complete")

	utils.PrintInfo("Updating configuration...")
	if err := cfg.Save(); err != nil {
		utils.PrintWarning("Warning: Rebuild succeeded but failed to save configuration: %v", err)
	}

	phpBuilder.CleanupSuccess()

	utils.PrintSuccess("Successfully rebuilt PHP %s with extensions: %s", version, strings.Join(extensions, ", "))
	return nil
}
