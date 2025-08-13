package php

import (
	"fmt"
	"strings"
	"time"

	"github.com/LumoSolutions/yerd/internal/builder"
	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/pkg/php"
	"github.com/fatih/color"
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
		color.New(color.FgRed).Printf("✗ Invalid PHP version: %s\n", phpVersion)
		return nil
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		color.New(color.FgRed).Printf("✗ Failed to load config: %v\n", err)
		return nil
	}

	if _, exists := cfg.InstalledPHP[phpVersion]; !exists {
		color.New(color.FgRed).Printf("✗ PHP %s is not installed. Use 'yerd php add %s' first\n", phpVersion, phpVersion)
		return nil
	}

	currentExtensions, exists := cfg.GetPHPExtensions(phpVersion)
	if !exists {
		color.New(color.FgRed).Printf("✗ No extension information found for PHP %s\n", phpVersion)
		return nil
	}

	if len(currentExtensions) == 0 {
		color.New(color.FgYellow).Printf("PHP %s has no extensions configured. Adding default extensions.\n", phpVersion)
		currentExtensions = []string{
			"mbstring", "bcmath", "opcache", "curl", "openssl", "zip",
			"sockets", "mysqli", "pdo-mysql", "gd", "jpeg", "freetype",
		}
		cfg.UpdatePHPExtensions(phpVersion, currentExtensions)
	}

	color.New(color.FgCyan, color.Bold).Printf("Rebuilding PHP %s with extensions: %s\n", phpVersion, strings.Join(currentExtensions, ", "))

	if err := forceRebuildPHP(cfg, phpVersion, currentExtensions); err != nil {
		return nil
	}

	return nil
}

// forceRebuildPHP performs the actual rebuild process with spinner animation.
// cfg: Configuration object, version: PHP version to rebuild, extensions: Extensions to include.
func forceRebuildPHP(cfg *config.Config, version string, extensions []string) error {
	color.New(color.FgYellow).Println("Force rebuilding PHP (no configuration backup needed)...")

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

	phpBuilder := builder.NewBuilder(version, extensions)
	err := phpBuilder.RebuildPHP()
	done <- true
	fmt.Print("\r")

	if err != nil {
		color.New(color.FgRed).Printf("✗ Failed to rebuild PHP %s: %v\n", version, err)
		phpBuilder.Cleanup()
		return fmt.Errorf("rebuild failed")
	}

	color.New(color.FgYellow).Println("Updating configuration...")
	if err := cfg.Save(); err != nil {
		color.New(color.FgRed).Printf("⚠️  Warning: Rebuild succeeded but failed to save configuration: %v\n", err)
	}

	phpBuilder.CleanupSuccess()

	color.New(color.FgGreen).Printf("✓ Successfully rebuilt PHP %s with extensions: %s\n", version, strings.Join(extensions, ", "))
	return nil
}
