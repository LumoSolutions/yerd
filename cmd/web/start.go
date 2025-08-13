package web

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/web"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start web services",
	Long: `Start nginx and dnsmasq services for local development.

This command starts both nginx and dnsmasq services:
  • nginx    - HTTP server and reverse proxy
  • dnsmasq  - DNS forwarder for local development

Example:
  yerd web start           # Start both nginx and dnsmasq`,
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintSplash()

		if !utils.CheckAndPromptForSudo("Web services management", "start") {
			return
		}

		manager, err := web.NewWebManager()
		if err != nil {
			utils.PrintError("Failed to create web manager: %v", err)
			return
		}

		fmt.Printf("Starting all web services:\n")
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

		missingServices := []string{}
		requiredServices := []string{"nginx", "dnsmasq"}
		for _, service := range requiredServices {
			found := false
			for _, installed := range installedServices {
				if installed == service {
					found = true
					break
				}
			}
			if !found {
				missingServices = append(missingServices, service)
			}
		}

		if len(missingServices) > 0 {
			utils.PrintError("Missing services: %v", missingServices)
			fmt.Printf("\n💡 To install missing services, run:\n")
			fmt.Printf("   sudo yerd web install\n")
			return
		}

		err = manager.StartAllServices()
		if err != nil {
			utils.PrintError("Failed to start services: %v", err)
			fmt.Printf("\n💡 Try running diagnostics:\n")
			fmt.Printf("   yerd doctor\n")
			return
		}

		utils.PrintSuccess("All web services started successfully")
	},
}