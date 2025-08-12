package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/pkg/php"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show YERD system status and conflicts",
	Long:  `Display current YERD configuration and check for potential conflicts with system PHP installations.`,
	Run:   runStatus,
}

type statusContext struct {
	cfg         *config.Config
	phpConflicts *utils.SystemPHPResult
	dirStatus   []utils.DirectoryStatus
	sysReq      *utils.SystemRequirementsResult
}

func runStatus(cmd *cobra.Command, args []string) {
	version.PrintSplash()
	
	ctx, err := initializeStatusContext()
	if err != nil {
		return
	}
	
	displayAllStatusSections(ctx)
}

func initializeStatusContext() (*statusContext, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		return nil, err
	}
	
	return &statusContext{
		cfg:          cfg,
		phpConflicts: utils.CheckSystemPHPConflicts(),
		dirStatus:    utils.GetYERDDirectoryStatus(),
		sysReq:       utils.CheckSystemRequirements(),
	}, nil
}

func displayAllStatusSections(ctx *statusContext) {
	displayYERDStatus(ctx.cfg)
	displaySystemPHPCheck(ctx.phpConflicts)
	displayDirectoryStatus(ctx.dirStatus)
	displayBuildEnvironment(ctx.sysReq)
	displayInstalledPHPVersions(ctx.cfg)
	displayPHPUpdateStatus(ctx.cfg)
}

func displayYERDStatus(cfg *config.Config) {
	fmt.Printf("📊 YERD Status\n")
	fmt.Printf("├─ Installed versions: %d\n", len(cfg.InstalledPHP))
	if cfg.CurrentCLI != "" {
		fmt.Printf("├─ Current CLI: PHP %s\n", cfg.CurrentCLI)
	} else {
		fmt.Printf("├─ Current CLI: None set\n")
	}
	fmt.Printf("└─ Config: ~/.config/yerd/config.json\n")
	fmt.Println()
}

func displaySystemPHPCheck(phpConflicts *utils.SystemPHPResult) {
	fmt.Printf("🔍 System PHP Check\n")
	if phpConflicts.Error != nil {
		fmt.Printf("├─ ❌ Error checking system PHP: %v\n", phpConflicts.Error)
	} else if phpConflicts.HasSystemPHP {
		displaySystemPHPConflict(phpConflicts)
	} else {
		fmt.Printf("└─ ✅ No conflicts - ready for YERD management\n")
	}
	fmt.Println()
}

func displaySystemPHPConflict(phpConflicts *utils.SystemPHPResult) {
	fmt.Printf("├─ ⚠️  System PHP detected\n")
	fmt.Printf("├─ Version: %s\n", phpConflicts.PHPInfo)
	fmt.Printf("├─ Type: %s\n", phpConflicts.PHPType)
	fmt.Printf("└─ Location: /usr/local/bin/php\n")
	fmt.Println()
	fmt.Printf("💡 Note: Remove system PHP to use YERD CLI versions\n")
}

func displayDirectoryStatus(dirStatus []utils.DirectoryStatus) {
	fmt.Printf("📁 Directory Status\n")
	for i, dir := range dirStatus {
		isLast := i == len(dirStatus)-1
		prefix := getTreePrefix(isLast)
		
		if dir.Exists {
			fmt.Printf("%s ✅ %s (%s)\n", prefix, dir.Path, dir.Description)
		} else {
			fmt.Printf("%s ❌ %s (%s) - missing\n", prefix, dir.Path, dir.Description)
		}
	}
	fmt.Printf("\n💡 Run with sudo to create missing directories\n")
	fmt.Println()
}

func displayBuildEnvironment(sysReq *utils.SystemRequirementsResult) {
	fmt.Printf("🔧 Build Environment\n")
	
	buildTools := []string{"gcc", "make", "wget", "tar"}
	for i, tool := range buildTools {
		isLast := i == len(buildTools)-1
		prefix := getTreePrefix(isLast)
		
		if available, exists := sysReq.BuildTools[tool]; exists && available {
			fmt.Printf("%s ✅ %s: Available\n", prefix, tool)
		} else {
			fmt.Printf("%s ❌ %s: Missing\n", prefix, tool)
		}
	}
	
	if !sysReq.AllAvailable {
		fmt.Printf("\n💡 Note: Missing build tools will be installed automatically during PHP installation\n")
	}
	fmt.Println()
}

func displayInstalledPHPVersions(cfg *config.Config) {
	fmt.Printf("📦 Installed PHP Versions\n")
	
	if len(cfg.InstalledPHP) == 0 {
		fmt.Printf("└─ No PHP versions installed\n")
		return
	}
	
	versionCount := len(cfg.InstalledPHP)
	currentIndex := 0
	
	for majorMinor, phpInfo := range cfg.InstalledPHP {
		currentIndex++
		isLast := currentIndex == versionCount
		displaySinglePHPVersion(cfg, majorMinor, phpInfo, isLast)
	}
	fmt.Println()
}

func displaySinglePHPVersion(cfg *config.Config, majorMinor string, phpInfo config.PHPInfo, isLast bool) {
	prefix := getTreePrefix(isLast)
	versionStatus := getPHPVersionStatus(cfg, majorMinor)
	
	fmt.Printf("%s %s\n", prefix, versionStatus)
	
	subPrefix := getSubTreePrefix(isLast)
	binaryPath := getPHPBinaryPath(majorMinor)
	iniPath := getPHPIniPath(majorMinor)
	
	fmt.Printf("%s├─ Binary: %s\n", subPrefix, binaryPath)
	fmt.Printf("%s├─ Config: %s\n", subPrefix, iniPath)
	fmt.Printf("%s└─ Install: %s\n", subPrefix, phpInfo.InstallPath)
	
	if !isLast {
		fmt.Printf("│\n")
	}
}

func displayPHPUpdateStatus(cfg *config.Config) {
	fmt.Printf("🔄 PHP Update Status\n")
	
	if len(cfg.InstalledPHP) == 0 {
		fmt.Printf("└─ No PHP versions to check\n")
		return
	}
	
	installedVersionsMap := buildInstalledVersionsMap(cfg)
	updateStatus, err := versions.CheckForUpdates(installedVersionsMap)
	
	if err != nil {
		fmt.Printf("└─ ❌ Could not check for updates: %v\n", err)
		return
	}
	
	displayUpdateResults(updateStatus)
}

func buildInstalledVersionsMap(cfg *config.Config) map[string]string {
	installedVersionsMap := make(map[string]string)
	for majorMinor, phpInfo := range cfg.InstalledPHP {
		binaryPath := php.GetBinaryPath(majorMinor)
		if output, err := utils.ExecuteCommand(binaryPath, "-v"); err == nil {
			installedVersionsMap[majorMinor] = output
		} else {
			installedVersionsMap[majorMinor] = phpInfo.Version
		}
	}
	return installedVersionsMap
}

func displayUpdateResults(updateStatus map[string]bool) {
	hasUpdates := false
	for majorMinor, hasUpdate := range updateStatus {
		if hasUpdate {
			fmt.Printf("├─ 🔄 PHP %s: Update available\n", majorMinor)
			hasUpdates = true
		} else {
			fmt.Printf("├─ ✅ PHP %s: Up to date\n", majorMinor)
		}
	}
	
	if hasUpdates {
		fmt.Printf("└─ 💡 Run 'yerd list' to see available updates\n")
	} else {
		fmt.Printf("└─ All installed PHP versions are up to date\n")
	}
}

func getTreePrefix(isLast bool) string {
	if isLast {
		return "└─"
	}
	return "├─"
}

func getSubTreePrefix(isLast bool) string {
	if isLast {
		return "   "
	}
	return "│  "
}

func getPHPVersionStatus(cfg *config.Config, majorMinor string) string {
	if cfg.CurrentCLI == majorMinor {
		return fmt.Sprintf("🎯 PHP %s (Current CLI)", majorMinor)
	}
	return fmt.Sprintf("📌 PHP %s", majorMinor)
}

func getPHPBinaryPath(majorMinor string) string {
	binaryPath, err := utils.GetPHPBinaryPath(majorMinor)
	if err != nil {
		return fmt.Sprintf("❌ %v", err)
	}
	return binaryPath
}

func getPHPIniPath(majorMinor string) string {
	iniPath, err := utils.GetPHPIniPath(majorMinor)
	if err != nil {
		return fmt.Sprintf("❌ %v", err)
	}
	return iniPath
}