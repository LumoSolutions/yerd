package php

import (
	"fmt"
	"os"

	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/pkg/php"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available and installed PHP versions",
	Long:  `Display a table of available PHP versions showing installation status and CLI configuration.`,
	Run:   runList,
}

type listData struct {
	cfg               *config.Config
	availableVersions []string
	latestVersions    map[string]string
	updateStatus      map[string]bool
	updateError       error
}

// runList displays a comprehensive table of available and installed PHP versions.
func runList(cmd *cobra.Command, args []string) {
	version.PrintSplash()

	data, err := loadListData()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	renderTable(data)
	printLegend(data.updateError)
}

// loadListData gathers all necessary information for the list command display.
// Returns listData struct containing config, versions, and update status or error if loading fails.
func loadListData() (*listData, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	data := &listData{
		cfg:               cfg,
		availableVersions: php.GetAvailableVersions(),
	}

	data.latestVersions, _, data.updateError = versions.GetLatestVersions()
	if data.updateError != nil {
		fmt.Printf("⚠️  Could not check for updates: %v\n", data.updateError)
	} else {
		data.updateStatus = getUpdateStatus(cfg)
	}

	return data, nil
}

// getUpdateStatus checks which installed PHP versions have updates available.
// cfg: Configuration object. Returns map of version strings to update availability booleans.
func getUpdateStatus(cfg *config.Config) map[string]bool {
	installedVersionsMap := make(map[string]string)
	for majorMinor, phpInfo := range cfg.InstalledPHP {
		binaryPath := php.GetBinaryPath(majorMinor)
		if output, err := utils.ExecuteCommand(binaryPath, "-v"); err == nil {
			installedVersionsMap[majorMinor] = output
		} else {
			installedVersionsMap[majorMinor] = phpInfo.Version
		}
	}

	updateStatus, _ := versions.CheckForUpdates(installedVersionsMap)
	return updateStatus
}

// renderTable creates and displays a formatted table of PHP version information.
// data: listData struct containing all display information.
func renderTable(data *listData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"VERSION", "INSTALLED", "CLI", "EXTENSIONS", "UPDATES"})

	for _, ver := range data.availableVersions {
		row := buildTableRow(ver, data)
		table.Append(row)
	}

	table.Render()
}

// buildTableRow constructs a single table row for a PHP version.
// ver: PHP version string, data: listData with all version info. Returns slice of row cell values.
func buildTableRow(ver string, data *listData) []string {
	installed := getInstallStatus(ver, data.cfg)
	cli := getCLIStatus(ver, data.cfg)
	extensions := getExtensionsStatus(ver, data.cfg)
	updates := getVersionUpdateStatus(ver, data)

	return []string{fmt.Sprintf("PHP %s", ver), installed, cli, extensions, updates}
}

// getInstallStatus returns installation status for a PHP version.
// ver: PHP version string, cfg: Configuration object. Returns "Yes" or "No".
func getInstallStatus(ver string, cfg *config.Config) string {
	if _, exists := cfg.InstalledPHP[ver]; exists {
		return "Yes"
	}
	return "No"
}

// getCLIStatus returns CLI indicator for a PHP version.
// ver: PHP version string, cfg: Configuration object. Returns "*" if current CLI, empty string otherwise.
func getCLIStatus(ver string, cfg *config.Config) string {
	if cfg.CurrentCLI == ver {
		return "*"
	}
	return ""
}

// getExtensionsStatus returns formatted extension count for an installed PHP version.
// ver: PHP version string, cfg: Configuration object. Returns extension count description or N/A.
func getExtensionsStatus(ver string, cfg *config.Config) string {
	if info, exists := cfg.InstalledPHP[ver]; exists {
		if len(info.Extensions) == 0 {
			return "None"
		}
		if len(info.Extensions) == 1 {
			return fmt.Sprintf("1: %s", info.Extensions[0])
		}
		return fmt.Sprintf("%d extensions", len(info.Extensions))
	}
	return "N/A"
}

// getVersionUpdateStatus returns update status based on whether version is installed.
// ver: PHP version string, data: listData with update information. Returns formatted update status.
func getVersionUpdateStatus(ver string, data *listData) string {
	if _, exists := data.cfg.InstalledPHP[ver]; exists {
		return getInstalledUpdateStatus(ver, data.updateStatus)
	}
	return getAvailableUpdateStatus(ver, data.latestVersions)
}

// getInstalledUpdateStatus returns update availability for installed PHP versions.
// ver: PHP version string, updateStatus: Map of update availability. Returns "Yes", "No", or "N/A".
func getInstalledUpdateStatus(ver string, updateStatus map[string]bool) string {
	if updateStatus == nil {
		return "N/A"
	}

	if hasUpdate, exists := updateStatus[ver]; exists && hasUpdate {
		return "Yes"
	}
	return "No"
}

// getAvailableUpdateStatus returns latest version info for available PHP versions.
// ver: PHP version string, latestVersions: Map of latest versions. Returns formatted latest version or N/A.
func getAvailableUpdateStatus(ver string, latestVersions map[string]string) string {
	if latestVersions == nil {
		return "N/A"
	}

	if latestVersion, exists := latestVersions[ver]; exists {
		return fmt.Sprintf("Latest: %s", latestVersion)
	}
	return "N/A"
}

// printLegend displays explanation of table symbols and handles network error cases.
// updateError: Error from update check, used to display appropriate legend information.
func printLegend(updateError error) {
	fmt.Println()
	fmt.Printf("Legend:\n")
	fmt.Printf("  * = Current CLI version\n")
	fmt.Printf("  Updates: Yes/No for installed versions, Latest: X.X.X for available versions\n")
	if updateError != nil {
		fmt.Printf("  N/A = Could not check for updates (network issue)\n")
	}
}
