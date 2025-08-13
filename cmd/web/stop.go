package web

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/web"
	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop web services",
	Long: `Stop nginx and dnsmasq services.

This command stops both nginx and dnsmasq services:
  • nginx    - HTTP server and reverse proxy
  • dnsmasq  - DNS forwarder for local development

Example:
  yerd web stop            # Stop both nginx and dnsmasq`,
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintSplash()

		if !utils.CheckAndPromptForSudo("Web services management", "stop") {
			return
		}

		manager, err := web.NewWebManager()
		if err != nil {
			utils.PrintError("Failed to create web manager: %v", err)
			return
		}

		fmt.Printf("Stopping all web services:\n")
		fmt.Printf("  • nginx    - HTTP server and reverse proxy\n")
		fmt.Printf("  • dnsmasq  - DNS forwarder for local development\n")
		fmt.Println()

		// Check if services are installed first
		installedServices := web.GetInstalledServices()
		if len(installedServices) == 0 {
			utils.PrintError("No web services are installed")
			fmt.Printf("\n💡 To install web services, run:\n")
			fmt.Printf("   sudo yerd web install\n")
			return
		}

		err = manager.StopAllServices()
		if err != nil {
			utils.PrintError("Failed to stop services: %v", err)
			fmt.Printf("\n💡 Some services may not be running or may have configuration issues\n")
			return
		}

		utils.PrintSuccess("All web services stopped successfully")
	},
}