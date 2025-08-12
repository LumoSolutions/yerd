package php

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/LumoSolutions/yerd/internal/config"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/versions"
	"github.com/LumoSolutions/yerd/pkg/php"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available and installed PHP versions",
	Long:  `Display a table of available PHP versions showing installation status and CLI configuration.`,
	Run:   runList,
}

type listData struct {
	cfg             *config.Config
	availableVersions []string
	latestVersions  map[string]string
	updateStatus    map[string]bool
	updateError     error
}

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

func renderTable(data *listData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"VERSION", "INSTALLED", "CLI", "UPDATES"})
	
	for _, ver := range data.availableVersions {
		row := buildTableRow(ver, data)
		table.Append(row)
	}
	
	table.Render()
}

func buildTableRow(ver string, data *listData) []string {
	installed := getInstallStatus(ver, data.cfg)
	cli := getCLIStatus(ver, data.cfg)
	updates := getVersionUpdateStatus(ver, data)
	
	return []string{fmt.Sprintf("PHP %s", ver), installed, cli, updates}
}

func getInstallStatus(ver string, cfg *config.Config) string {
	if _, exists := cfg.InstalledPHP[ver]; exists {
		return "Yes"
	}
	return "No"
}

func getCLIStatus(ver string, cfg *config.Config) string {
	if cfg.CurrentCLI == ver {
		return "*"
	}
	return ""
}

func getVersionUpdateStatus(ver string, data *listData) string {
	if _, exists := data.cfg.InstalledPHP[ver]; exists {
		return getInstalledUpdateStatus(ver, data.updateStatus)
	}
	return getAvailableUpdateStatus(ver, data.latestVersions)
}

func getInstalledUpdateStatus(ver string, updateStatus map[string]bool) string {
	if updateStatus == nil {
		return "N/A"
	}
	
	if hasUpdate, exists := updateStatus[ver]; exists && hasUpdate {
		return "Yes"
	}
	return "No"
}

func getAvailableUpdateStatus(ver string, latestVersions map[string]string) string {
	if latestVersions == nil {
		return "N/A"
	}
	
	if latestVersion, exists := latestVersions[ver]; exists {
		return fmt.Sprintf("Latest: %s", latestVersion)
	}
	return "N/A"
}

func printLegend(updateError error) {
	fmt.Println()
	fmt.Printf("Legend:\n")
	fmt.Printf("  * = Current CLI version\n")
	fmt.Printf("  Updates: Yes/No for installed versions, Latest: X.X.X for available versions\n")
	if updateError != nil {
		fmt.Printf("  N/A = Could not check for updates (network issue)\n")
	}
}