package php

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/installer"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/pkg/php"
)

var UpdateCmd = &cobra.Command{
	Use:   "update [version]",
	Short: "Update installed PHP versions to latest releases",
	Long: `Check for and install updates to PHP versions. This command ignores the local cache 
and checks php.net directly for the latest versions.

Examples:
  yerd php update           # Update all installed PHP versions that have updates
  yerd php update 8.4       # Update only PHP 8.4 if an update is available
  yerd php update -y        # Update all versions without confirmation
  yerd php update 8.4 -y    # Update PHP 8.4 without confirmation`,
	Args: cobra.MaximumNArgs(1),
	Run:  runUpdate,
}

var forceYes bool

func init() {
	UpdateCmd.Flags().BoolVarP(&forceYes, "yes", "y", false, "Automatically confirm updates without prompting")
}

type updateContext struct {
	cfg              *config.Config
	targetVersion    string
	installedVersions map[string]string
	updates          map[string]bool
	availableUpdates map[string]string
	successful       int
	failed           int
}

func runUpdate(cmd *cobra.Command, args []string) {
	version.PrintSplash()
	
	ctx, err := initializeUpdate(args)
	if err != nil {
		return
	}
	
	if !checkAndDisplayUpdates(ctx) {
		return
	}
	
	if !confirmUpdates(ctx) {
		return
	}
	
	performUpdates(ctx)
	displayUpdateSummary(ctx)
}

func initializeUpdate(args []string) (*updateContext, error) {
	if err := utils.CheckInstallPermissions(); err != nil {
		fmt.Printf("❌ Update requires elevated permissions: %v\n\nTry running: sudo yerd update", err)
		if len(args) > 0 {
			fmt.Printf(" %s", args[0])
		}
		fmt.Println()
		return nil, err
	}
	
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		return nil, err
	}
	
	if len(cfg.InstalledPHP) == 0 {
		fmt.Printf("📦 No PHP versions installed\n")
		fmt.Printf("💡 Install a PHP version first: sudo yerd php add 8.4\n")
		return nil, fmt.Errorf("no PHP versions installed")
	}
	
	ctx := &updateContext{cfg: cfg}
	
	if err := ctx.parseTargetVersion(args); err != nil {
		return nil, err
	}
	
	ctx.buildInstalledVersionsMap()
	return ctx, nil
}

func (ctx *updateContext) parseTargetVersion(args []string) error {
	if len(args) == 0 {
		return nil
	}
	
	ctx.targetVersion = php.FormatVersion(args[0])
	if !php.IsValidVersion(ctx.targetVersion) {
		fmt.Printf("❌ Invalid PHP version: %s\n", args[0])
		fmt.Printf("💡 Valid versions: %s\n", strings.Join(php.GetAvailableVersions(), ", "))
		return fmt.Errorf("invalid version")
	}
	
	if _, exists := ctx.cfg.InstalledPHP[ctx.targetVersion]; !exists {
		fmt.Printf("❌ PHP %s is not installed\n", ctx.targetVersion)
		fmt.Printf("💡 Install it first: sudo yerd php add %s\n", ctx.targetVersion)
		return fmt.Errorf("version not installed")
	}
	
	return nil
}

func (ctx *updateContext) buildInstalledVersionsMap() {
	ctx.installedVersions = make(map[string]string)
	for majorMinor, phpInfo := range ctx.cfg.InstalledPHP {
		if ctx.targetVersion != "" && majorMinor != ctx.targetVersion {
			continue
		}
		
		binaryPath := php.GetBinaryPath(majorMinor)
		if output, err := utils.ExecuteCommand(binaryPath, "-v"); err == nil {
			ctx.installedVersions[majorMinor] = output
		} else {
			ctx.installedVersions[majorMinor] = phpInfo.Version
		}
	}
}

func checkAndDisplayUpdates(ctx *updateContext) bool {
	if len(ctx.installedVersions) == 0 {
		fmt.Printf("❌ No matching PHP versions found to update\n")
		return false
	}
	
	spinner := utils.NewLoadingSpinner("🔍 Checking for updates from php.net...")
	spinner.Start()
	
	updates, availableUpdates, err := versions.CheckForUpdatesFresh(ctx.installedVersions)
	spinner.Stop("")
	
	if err != nil {
		fmt.Printf("❌ Failed to check for updates: %v\n", err)
		fmt.Printf("💡 Check your internet connection and try again\n")
		return false
	}
	
	ctx.updates = updates
	ctx.availableUpdates = availableUpdates
	
	updateCount := countUpdates(updates)
	if updateCount == 0 {
		printNoUpdatesMessage(ctx.targetVersion)
		return false
	}
	
	displayAvailableUpdates(ctx, updateCount)
	return true
}

func countUpdates(updates map[string]bool) int {
	count := 0
	for _, hasUpdate := range updates {
		if hasUpdate {
			count++
		}
	}
	return count
}

func printNoUpdatesMessage(targetVersion string) {
	if targetVersion != "" {
		fmt.Printf("✅ PHP %s is already up to date\n", targetVersion)
	} else {
		fmt.Printf("✅ All installed PHP versions are up to date\n")
	}
}

func displayAvailableUpdates(ctx *updateContext, updateCount int) {
	fmt.Printf("🔄 Found %d update(s) available:\n", updateCount)
	for majorMinor, hasUpdate := range ctx.updates {
		if hasUpdate {
			currentVersion := versions.ExtractVersionFromString(ctx.installedVersions[majorMinor])
			newVersion := ctx.availableUpdates[majorMinor]
			fmt.Printf("├─ PHP %s: %s → %s\n", majorMinor, currentVersion, newVersion)
		}
	}
	fmt.Println()
}

func confirmUpdates(ctx *updateContext) bool {
	updateCount := countUpdates(ctx.updates)
	
	if forceYes {
		printAutoUpdateMessage(ctx.targetVersion, updateCount)
		return true
	}
	
	return promptUserConfirmation(ctx.targetVersion, updateCount)
}

func printAutoUpdateMessage(targetVersion string, updateCount int) {
	if targetVersion != "" {
		fmt.Printf("🔄 Auto-updating PHP %s...\n", targetVersion)
	} else {
		fmt.Printf("🔄 Auto-updating all %d PHP version(s)...\n", updateCount)
	}
}

func promptUserConfirmation(targetVersion string, updateCount int) bool {
	if targetVersion != "" {
		fmt.Printf("🔄 Update PHP %s? (y/N): ", targetVersion)
	} else {
		fmt.Printf("🔄 Update all %d PHP version(s)? (y/N): ", updateCount)
	}
	
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	
	if response != "y" && response != "yes" {
		fmt.Printf("❌ Update cancelled\n")
		return false
	}
	return true
}

func performUpdates(ctx *updateContext) {
	for majorMinor, hasUpdate := range ctx.updates {
		if hasUpdate {
			updateSinglePHPVersion(ctx, majorMinor)
		}
	}
}

func updateSinglePHPVersion(ctx *updateContext, majorMinor string) {
	fmt.Printf("\n🔄 Updating PHP %s...\n", majorMinor)
	
	wasCLI := ctx.cfg.CurrentCLI == majorMinor
	originalConfig := ctx.cfg.InstalledPHP[majorMinor]
	
	prepareForUpdate(ctx, majorMinor, wasCLI)
	
	if err := installer.InstallPHP(majorMinor, false); err != nil {
		rollbackUpdate(ctx, majorMinor, originalConfig, wasCLI, err)
		return
	}
	
	finishUpdate(ctx, majorMinor, originalConfig, wasCLI)
}

func prepareForUpdate(ctx *updateContext, majorMinor string, wasCLI bool) {
	fmt.Printf("📦 Installing updated PHP %s...\n", majorMinor)
	delete(ctx.cfg.InstalledPHP, majorMinor)
	if wasCLI {
		ctx.cfg.CurrentCLI = ""
	}
}

func rollbackUpdate(ctx *updateContext, majorMinor string, originalConfig config.PHPInfo, wasCLI bool, err error) {
	ctx.cfg.InstalledPHP[majorMinor] = originalConfig
	if wasCLI {
		ctx.cfg.CurrentCLI = majorMinor
	}
	ctx.cfg.Save()
	
	fmt.Printf("❌ Failed to install updated PHP %s: %v\n", majorMinor, err)
	fmt.Printf("💡 Your existing PHP %s installation is still available\n", majorMinor)
	ctx.failed++
}

func finishUpdate(ctx *updateContext, majorMinor string, originalConfig config.PHPInfo, wasCLI bool) {
	fmt.Printf("🗑️  Cleaning up old PHP %s installation...\n", majorMinor)
	if err := cleanupOldVersion(majorMinor, originalConfig); err != nil {
		fmt.Printf("⚠️  Warning: Failed to cleanup old files for PHP %s: %v\n", majorMinor, err)
		fmt.Printf("💡 New PHP %s is installed and working, old files may remain\n", majorMinor)
	}
	
	if wasCLI {
		if err := restoreCLIVersion(ctx.cfg, majorMinor); err != nil {
			fmt.Printf("⚠️  Warning: %v\n", err)
			fmt.Printf("💡 Run 'sudo yerd php cli %s' to fix this\n", majorMinor)
		}
	}
	
	fmt.Printf("✅ PHP %s updated successfully to %s\n", majorMinor, ctx.availableUpdates[majorMinor])
	ctx.successful++
}

func displayUpdateSummary(ctx *updateContext) {
	fmt.Printf("\n📊 Update Summary:\n")
	if ctx.successful > 0 {
		fmt.Printf("├─ ✅ Successfully updated: %d\n", ctx.successful)
	}
	if ctx.failed > 0 {
		fmt.Printf("├─ ❌ Failed to update: %d\n", ctx.failed)
	}
	fmt.Printf("└─ 🔄 Total processed: %d\n", ctx.successful+ctx.failed)
	
	if ctx.successful > 0 {
		printVerificationInstructions(ctx.updates)
	}
}

func printVerificationInstructions(updates map[string]bool) {
	fmt.Printf("\n💡 Verify your updates:\n")
	for majorMinor, hasUpdate := range updates {
		if hasUpdate {
			fmt.Printf("   php%s -v\n", majorMinor)
		}
	}
}

func restoreCLIVersion(cfg *config.Config, majorMinor string) error {
	cfg.CurrentCLI = majorMinor
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}
	
	binaryPath := php.GetBinaryPath(majorMinor)
	cliPath := utils.SystemBinDir + "/php"
	if err := utils.CreateSymlink(binaryPath, cliPath); err != nil {
		return fmt.Errorf("failed to restore CLI symlink: %v", err)
	}
	
	return nil
}

func cleanupOldVersion(version string, originalConfig config.PHPInfo) error {
	installPath := originalConfig.InstallPath
	if installPath != "" && utils.FileExists(installPath) {
		output, err := utils.ExecuteCommand("rm", "-rf", installPath)
		if err != nil {
			return fmt.Errorf("failed to remove old installation directory %s: %v (output: %s)", installPath, err, output)
		}
	}
	
	return nil
}