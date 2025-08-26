package php

import (
	"fmt"
	"os"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	phpinstaller "github.com/lumosolutions/yerd/internal/installers/php"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func BuildListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists installed PHP versions",
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()

			versions := constants.GetAvailablePhpVersions()
			latestVersions, _, _ := phpinstaller.GetLatestVersions()

			rows := [][]string{}
			headers := []string{"VERSION", "INSTALLED", "CLI", "EXTENSIONS", "UPDATES"}

			for _, version := range versions {
				if data, installed := config.GetInstalledPhpInfo(version); installed {
					rows = append(rows, []string{
						data.Version,
						data.InstalledVersion,
						friendlyBool(data.IsCLI),
						fmt.Sprintf("%d", len(data.Extensions)),
						friendlyBool(data.InstalledVersion != latestVersions[version]),
					})
				}
			}

			if len(rows) == 0 {
				fmt.Println("No YERD PHP versions installed")
				fmt.Println("Run 'sudo yerd php {version} install' to get started")
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.Header(headers)
			table.Bulk(rows)

			table.Render()
		},
	}

	return cmd
}

func friendlyBool(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
}
