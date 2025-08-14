package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/internal/web"
	"github.com/LumoSolutions/yerd/pkg/php"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show YERD system status and conflicts",
	Long:  `Display current YERD configuration and check for potential conflicts with system PHP installations.`,
	Run:   runStatus,
}

type statusContext struct {
	cfg          *config.Config
	phpConflicts *utils.SystemPHPResult
	dirStatus    []utils.DirectoryStatus
	sysReq       *utils.SystemRequirementsResult
}

// runStatus executes the status command, displaying comprehensive system information.
func runStatus(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	ctx, err := initializeStatusContext()
	if err != nil {
		return
	}

	displayAllStatusSections(ctx)
}

// initializeStatusContext loads configuration and gathers system information for status display.
// Returns a statusContext struct containing all necessary data or an error if config loading fails.
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

// displayAllStatusSections renders all status information sections in order.
// Takes a statusContext containing pre-gathered system information.
func displayAllStatusSections(ctx *statusContext) {
	displayYERDStatus(ctx.cfg)
	displaySystemPHPCheck(ctx.phpConflicts)
	displayDirectoryStatus(ctx.dirStatus)
	displayBuildEnvironment(ctx.sysReq)
	displayInstalledPHPVersions(ctx.cfg)
	displayWebServicesStatus()
	displayPHPUpdateStatus(ctx.cfg)
}

// displayYERDStatus shows basic YERD installation information including installed versions and current CLI.
// cfg: Configuration object containing installed PHP versions and CLI settings.
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

// displaySystemPHPCheck shows system PHP conflict detection results.
// phpConflicts: Result of system PHP detection including conflict status and version info.
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

// displaySystemPHPConflict shows detailed information when system PHP conflicts are detected.
// phpConflicts: SystemPHPResult containing details about the conflicting PHP installation.
func displaySystemPHPConflict(phpConflicts *utils.SystemPHPResult) {
	fmt.Printf("├─ ⚠️  System PHP detected\n")
	fmt.Printf("├─ Version: %s\n", phpConflicts.PHPInfo)
	fmt.Printf("├─ Type: %s\n", phpConflicts.PHPType)
	fmt.Printf("└─ Location: /usr/local/bin/php\n")
	fmt.Println()
	fmt.Printf("💡 Note: Remove system PHP to use YERD CLI versions\n")
}

// displayDirectoryStatus shows the existence status of YERD directories.
// dirStatus: Slice of DirectoryStatus objects indicating which directories exist and their purposes.
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

// displayBuildEnvironment shows availability of required build tools for PHP compilation.
// sysReq: SystemRequirementsResult containing build tool availability status.
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

// displayInstalledPHPVersions shows detailed information about all installed PHP versions.
// cfg: Configuration object containing installed PHP version details.
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

// displaySinglePHPVersion renders information for one PHP installation.
// cfg: Configuration object, majorMinor: PHP version string, phpInfo: Installation details, isLast: Controls tree formatting.
func displaySinglePHPVersion(cfg *config.Config, majorMinor string, phpInfo config.PHPInfo, isLast bool) {
	prefix := getTreePrefix(isLast)
	versionStatus := getPHPVersionStatus(cfg, majorMinor)

	fmt.Printf("%s %s\n", prefix, versionStatus)

	subPrefix := getSubTreePrefix(isLast)
	binaryPath := getPHPBinaryPath(majorMinor)
	iniPath := getPHPIniPath(majorMinor)
	fpmSockPath := getFPMSocketPath(majorMinor)
	fpmPoolConfig := getFPMPoolConfigPath(majorMinor)

	fmt.Printf("%s├─ Binary: %s\n", subPrefix, binaryPath)
	fmt.Printf("%s├─ Config: %s\n", subPrefix, iniPath)
	fmt.Printf("%s├─ FPM Socket: %s\n", subPrefix, fpmSockPath)
	fmt.Printf("%s├─ FPM Pool: %s\n", subPrefix, fpmPoolConfig)
	fmt.Printf("%s└─ Install: %s\n", subPrefix, phpInfo.InstallPath)

	if !isLast {
		fmt.Printf("│\n")
	}
}

// displayPHPUpdateStatus checks and displays update availability for installed PHP versions.
// cfg: Configuration object containing installed PHP version information.
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

// buildInstalledVersionsMap creates a map of PHP versions to their actual version strings.
// cfg: Configuration object. Returns map where keys are major.minor versions and values are full version strings.
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

// displayUpdateResults shows which PHP versions have updates available.
// updateStatus: Map where keys are PHP versions and values indicate if updates are available.
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

// getTreePrefix returns appropriate tree drawing characters for list formatting.
// isLast: If true, returns characters for final item; otherwise returns characters for middle items.
func getTreePrefix(isLast bool) string {
	if isLast {
		return "└─"
	}
	return "├─"
}

// getSubTreePrefix returns appropriate indentation for sub-items in tree display.
// isLast: If true, returns spacing for final parent item; otherwise returns continued tree line.
func getSubTreePrefix(isLast bool) string {
	if isLast {
		return "   "
	}
	return "│  "
}

// getPHPVersionStatus returns a formatted status string for a PHP version.
// cfg: Configuration object, majorMinor: PHP version to check. Returns string with emoji and CLI indicator.
func getPHPVersionStatus(cfg *config.Config, majorMinor string) string {
	if cfg.CurrentCLI == majorMinor {
		return fmt.Sprintf("🎯 PHP %s (Current CLI)", majorMinor)
	}
	return fmt.Sprintf("📌 PHP %s", majorMinor)
}

// getPHPBinaryPath returns the path to a PHP version's binary or an error message.
// majorMinor: PHP version string. Returns formatted binary path or error description.
func getPHPBinaryPath(majorMinor string) string {
	binaryPath, err := utils.GetPHPBinaryPath(majorMinor)
	if err != nil {
		return fmt.Sprintf("❌ %v", err)
	}
	return binaryPath
}

// getPHPIniPath returns the path to a PHP version's ini file or an error message.
// majorMinor: PHP version string. Returns formatted ini path or error description.
func getPHPIniPath(majorMinor string) string {
	iniPath, err := utils.GetPHPIniPath(majorMinor)
	if err != nil {
		return fmt.Sprintf("❌ %v", err)
	}
	return iniPath
}

// displayWebServicesStatus shows information about installed web services (nginx)
func displayWebServicesStatus() {
	services := []string{"nginx"}
	installedServices := make([]string, 0)

	for _, service := range services {
		if web.IsServiceInstalled(service) {
			installedServices = append(installedServices, service)
		}
	}

	if len(installedServices) == 0 {
		return
	}

	fmt.Printf("🌐 Web Services\n")

	for i, service := range installedServices {
		isLast := i == len(installedServices)-1
		displaySingleWebService(service, isLast)
	}
	fmt.Println()
}

// displaySingleWebService shows detailed information for one web service
func displaySingleWebService(service string, isLast bool) {
	prefix := getTreePrefix(isLast)
	config, _ := web.GetServiceConfig(service)
	
	serviceStatus := getWebServiceStatus(service)
	fmt.Printf("%s %s %s %s\n", prefix, serviceStatus, service, config.Version)

	subPrefix := getSubTreePrefix(isLast)
	binaryPath := web.GetServiceBinaryPath(service)
	configPath := getWebServiceConfigPath(service)
	runningStatus := getWebServiceRunningStatus(service)

	fmt.Printf("%s├─ Binary: %s\n", subPrefix, binaryPath)
	fmt.Printf("%s├─ Config: %s\n", subPrefix, configPath)
	fmt.Printf("%s└─ Status: %s\n", subPrefix, runningStatus)

	if !isLast {
		fmt.Printf("│\n")
	}
}

// getWebServiceStatus returns a formatted status string for a web service
func getWebServiceStatus(service string) string {
	return "🔧"
}

// getWebServiceConfigPath returns the configuration file path for a web service
func getWebServiceConfigPath(service string) string {
	switch service {
	case "nginx":
		return filepath.Join(web.GetServiceConfigPath(service), "nginx.conf")
	default:
		return web.GetServiceConfigPath(service)
	}
}

// getWebServiceRunningStatus checks if a web service is currently running
func getWebServiceRunningStatus(service string) string {
	switch service {
	case "nginx":
		if isNginxRunning() {
			return "🟢 Running"
		}
		return "🔴 Stopped"
	default:
		return "❓ Unknown"
	}
}

// isNginxRunning checks if nginx is currently running by checking for the process
func isNginxRunning() bool {
	binaryPath := web.GetServiceBinaryPath("nginx")
	pidPath := filepath.Join(web.GetServiceRunPath("nginx"), "nginx.pid")
	
	if utils.FileExists(pidPath) {
		if _, err := utils.ExecuteCommand("pgrep", "-f", binaryPath); err == nil {
			return true
		}
	}
	
	if _, err := utils.ExecuteCommand("pgrep", "-f", binaryPath); err == nil {
		return true
	}
	
	return false
}

// getFPMSocketPath returns the FPM socket path for a PHP version
func getFPMSocketPath(majorMinor string) string {
	sockPath := filepath.Join(utils.FPMSockDir, fmt.Sprintf("php%s-fpm.sock", majorMinor))
	if utils.FileExists(sockPath) {
		return sockPath
	}
	return fmt.Sprintf("❌ %s (not found)", sockPath)
}

// getFPMPoolConfigPath returns the FPM pool configuration path for a PHP version
func getFPMPoolConfigPath(majorMinor string) string {
	configDir := filepath.Join(utils.YerdEtcDir, "php"+majorMinor)
	poolConfigPath := filepath.Join(configDir, "php-fpm.d", "www.conf")
	if utils.FileExists(poolConfigPath) {
		return poolConfigPath
	}
	return fmt.Sprintf("❌ %s (not found)", poolConfigPath)
}

