package php

import (
	"fmt"
	"sort"
	"time"

	"github.com/LumoSolutions/yerd/internal/builder"
	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/pkg/extensions"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ExtensionsCmd = &cobra.Command{
	Use:   "extensions <version>",
	Short: "Show installed and available PHP extensions",
	Long: `Show installed and available PHP extensions for a specific version.

Examples:
  yerd php extensions 8.3                       # Show installed and available extensions
  yerd php extensions add 8.3 mysqli zip        # Add extensions
  yerd php extensions remove 8.3 gd curl        # Remove extensions  
  yerd php extensions replace 8.3 mysqli zip    # Replace all extensions with specified ones`,
	Args: cobra.ExactArgs(1),
	RunE: runExtensions,
}

var AddExtensionsCmd = &cobra.Command{
	Use:   "add <version> <extension1> [extension2...]",
	Short: "Add PHP extensions to a version",
	Long:  `Add one or more PHP extensions to the specified version.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAddExtensions,
}

var RemoveExtensionsCmd = &cobra.Command{
	Use:   "remove <version> <extension1> [extension2...]",
	Short: "Remove PHP extensions from a version",
	Long:  `Remove one or more PHP extensions from the specified version.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runRemoveExtensions,
}

var ReplaceExtensionsCmd = &cobra.Command{
	Use:   "replace <version> <extension1> [extension2...]",
	Short: "Replace all extensions with the specified ones",
	Long:  `Remove all existing extensions and replace them with the specified ones.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runReplaceExtensions,
}

func init() {
	ExtensionsCmd.AddCommand(AddExtensionsCmd)
	ExtensionsCmd.AddCommand(RemoveExtensionsCmd)
	ExtensionsCmd.AddCommand(ReplaceExtensionsCmd)
}

// runExtensions displays installed and available PHP extensions for a specific version.
func runExtensions(cmd *cobra.Command, args []string) error {
	version.PrintSplash()

	phpVersion, cfg, err := ValidatePHPVersionConfigAndInstallation(args[0])
	if err != nil {
		return err
	}

	return listExtensions(cfg, phpVersion)
}

// runAddExtensions adds one or more PHP extensions to an installed version with validation.
func runAddExtensions(cmd *cobra.Command, args []string) error {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo(OpExtensionManagement, "add", args[0]) {
		return nil
	}

	phpVersion, cfg, err := ValidatePHPVersionConfigAndInstallation(args[0])
	if err != nil {
		return err
	}

	extensionsToAdd := args[1:]
	valid, invalid := extensions.ValidateExtensions(extensionsToAdd)
	
	if err := utils.PrintInvalidExtensionsWithSuggestions(invalid); err != nil {
		return err
	}

	currentExtensions, _ := cfg.GetPHPExtensions(phpVersion)
	currentMap := utils.CreateExtensionMap(currentExtensions)

	var newExtensions []string
	var alreadyInstalled []string

	for _, ext := range valid {
		if currentMap[ext] {
			alreadyInstalled = append(alreadyInstalled, ext)
		} else {
			currentMap[ext] = true
			newExtensions = append(newExtensions, ext)
		}
	}

	utils.PrintExtensionList("Already installed", alreadyInstalled)

	if len(newExtensions) == 0 {
		utils.PrintWarning("No new extensions to add for PHP %s", phpVersion)
		return nil
	}

	utils.PrintInfo("Adding extensions to PHP %s:", phpVersion)
	sort.Strings(newExtensions)
	utils.PrintExtensionsGrid(newExtensions)
	fmt.Println()

	finalExtensions := utils.SliceFromExtensionMap(currentMap)
	return applyExtensionChanges(cfg, phpVersion, finalExtensions)
}

// runRemoveExtensions removes one or more PHP extensions from an installed version.
func runRemoveExtensions(cmd *cobra.Command, args []string) error {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo(OpExtensionManagement, "remove", args[0]) {
		return nil
	}

	phpVersion, cfg, err := ValidatePHPVersionConfigAndInstallation(args[0])
	if err != nil {
		return err
	}

	extensionsToRemove := args[1:]
	currentExtensions, _ := cfg.GetPHPExtensions(phpVersion)
	currentMap := utils.CreateExtensionMap(currentExtensions)

	var removedExtensions []string
	var notInstalled []string

	for _, ext := range extensionsToRemove {
		if currentMap[ext] {
			delete(currentMap, ext)
			removedExtensions = append(removedExtensions, ext)
		} else {
			notInstalled = append(notInstalled, ext)
		}
	}

	utils.PrintExtensionList("Not installed", notInstalled)

	if len(removedExtensions) == 0 {
		utils.PrintWarning("No extensions to remove from PHP %s", phpVersion)
		return nil
	}

	utils.PrintInfo("Removing extensions from PHP %s:", phpVersion)
	sort.Strings(removedExtensions)
	utils.PrintExtensionsGrid(removedExtensions)
	fmt.Println()

	finalExtensions := utils.SliceFromExtensionMap(currentMap)
	return applyExtensionChanges(cfg, phpVersion, finalExtensions)
}

// runReplaceExtensions replaces all existing extensions with a new set of extensions.
func runReplaceExtensions(cmd *cobra.Command, args []string) error {
	version.PrintSplash()

	if !utils.CheckAndPromptForSudo(OpExtensionManagement, "replace", args[0]) {
		return nil
	}

	phpVersion, cfg, err := ValidatePHPVersionConfigAndInstallation(args[0])
	if err != nil {
		return err
	}

	newExtensions := args[1:]
	valid, invalid := extensions.ValidateExtensions(newExtensions)
	
	if err := utils.PrintInvalidExtensionsWithSuggestions(invalid); err != nil {
		return err
	}

	currentExtensions, _ := cfg.GetPHPExtensions(phpVersion)
	sort.Strings(valid)
	sort.Strings(currentExtensions)

	if fmt.Sprintf("%v", valid) == fmt.Sprintf("%v", currentExtensions) {
		utils.PrintWarning("Extensions for PHP %s are already set to:", phpVersion)
		utils.PrintExtensionsGrid(valid)
		return nil
	}

	utils.PrintInfo("Replacing all extensions for PHP %s with:", phpVersion)
	utils.PrintExtensionsGrid(valid)
	fmt.Println()

	return applyExtensionChanges(cfg, phpVersion, valid)
}

// listExtensions displays installed and available extensions in grid format for a PHP version.
// cfg: Configuration object, version: PHP version string. Returns error if version data not found.
func listExtensions(cfg *config.Config, version string) error {
	installedExtensions, exists := cfg.GetPHPExtensions(version)
	if !exists {
		return fmt.Errorf(ErrNoExtensionInfo, version)
	}

	utils.PrintInfo("PHP %s Extensions:", version)
	fmt.Println()

	utils.PrintSuccess("INSTALLED:")
	if len(installedExtensions) == 0 {
		fmt.Println("  No extensions installed")
	} else {
		sort.Strings(installedExtensions)
		utils.PrintExtensionsGrid(installedExtensions)
	}

	fmt.Println()

	utils.PrintInfo("AVAILABLE:")
	installedMap := make(map[string]bool)
	for _, ext := range installedExtensions {
		installedMap[ext] = true
	}

	var availableExtensions []string
	for name := range extensions.AvailableExtensions {
		if !installedMap[name] {
			availableExtensions = append(availableExtensions, name)
		}
	}

	if len(availableExtensions) == 0 {
		fmt.Println("  All available extensions are already installed")
	} else {
		sort.Strings(availableExtensions)
		utils.PrintExtensionsGrid(availableExtensions)
	}

	fmt.Println()
	utils.PrintInfo("USAGE:")
	fmt.Printf("  yerd php extensions add %s <extension>     # Add extensions\n", version)
	fmt.Printf("  yerd php extensions remove %s <extension>  # Remove extensions\n", version)
	fmt.Printf("  yerd php extensions replace %s <extension> # Replace all extensions\n", version)
	fmt.Printf("  yerd php rebuild %s                        # Force rebuild with current extensions\n", version)

	return nil
}

// applyExtensionChanges updates PHP extension configuration and triggers rebuild with rollback support.
// cfg: Configuration object, version: PHP version, extensions: New extension list. Returns error if rebuild fails.
func applyExtensionChanges(cfg *config.Config, version string, extensions []string) error {
	snapshot := cfg.CreateSnapshot(version)
	color.New(color.FgBlue).Printf("📸 Created configuration backup for PHP %s\n", version)

	cfg.UpdatePHPExtensions(version, extensions)

	color.New(color.FgYellow).Println("Rebuilding PHP with new extensions...")

	spinner := []string{"|", "/", "-", "\\"}
	done := make(chan bool)

	go func() {
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Building PHP %s with extensions... ", spinner[i%len(spinner)], version)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()

	err := rebuildPHPWithExtensions(version, extensions)
	done <- true
	fmt.Print("\r\033[K") // Clear the entire line

	if err != nil {
		color.New(color.FgRed).Printf("✗ Failed to rebuild PHP %s: %v\n", version, err)
		cfg.RestoreSnapshot(snapshot)
		color.New(color.FgYellow).Printf("↺ Restored configuration to previous state\n")

		if saveErr := cfg.Save(); saveErr != nil {
			color.New(color.FgRed).Printf("⚠️  Warning: Failed to save restored configuration: %v\n", saveErr)
		} else {
			color.New(color.FgGreen).Println("✓ Configuration restored successfully")
		}

		return fmt.Errorf(ErrRebuildFailed)
	}

	color.New(color.FgYellow).Println("Saving new configuration...")
	if err := cfg.Save(); err != nil {
		color.New(color.FgRed).Printf("⚠️  Warning: Build succeeded but failed to save configuration: %v\n", err)
	}

	color.New(color.FgGreen).Printf("✓ Successfully updated PHP %s extensions\n", version)
	return nil
}

// rebuildPHPWithExtensions performs PHP rebuild with specified extensions using the builder.
// version: PHP version string, extensions: Extension list. Returns error if build fails.
func rebuildPHPWithExtensions(version string, extensions []string) error {
	phpBuilder, err := builder.NewBuilder(version, extensions)
	if err != nil {
		return fmt.Errorf(ErrFailedToCreateBuilder, err)
	}
	defer phpBuilder.Cleanup()

	return phpBuilder.RebuildPHP()
}
