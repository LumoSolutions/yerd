package web

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/internal/web"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install web services",
	Long:  `Install nginx and dnsmasq for local development.

Examples:
  yerd web install     # Install if not already installed
  yerd web install -f  # Force reinstall even if already installed`,
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintSplash()

		if !utils.CheckAndPromptForSudo("Web services installation", "install") {
			return
		}

		force, _ := cmd.Flags().GetBool("force")
		services := []string{"nginx", "dnsmasq"}

		if force {
			fmt.Printf("Force reinstalling web services:\n")
		} else {
			fmt.Printf("Installing web services:\n")
		}
		fmt.Printf("  • nginx 1.29.1   - High-performance HTTP server and reverse proxy\n")
		fmt.Printf("  • dnsmasq 2.91   - Lightweight DNS forwarder and DHCP server\n")
		fmt.Println()

		var failed []string

		for _, service := range services {
			// Check if service is already installed (skip check if force flag is set)
			if !force && web.IsServiceInstalled(service) {
				utils.PrintWarning("Service '%s' is already installed, skipping", service)
				utils.PrintInfo("Use -f or --force to reinstall")
				continue
			}

			// Create installer
			installer, err := web.NewWebInstaller(service)
			if err != nil {
				utils.PrintError("Failed to create installer for %s: %v", service, err)
				failed = append(failed, service)
				continue
			}

			// If force flag is set and service is installed, note it for replacement after build
			var wasInstalled bool
			if force && web.IsServiceInstalled(service) {
				wasInstalled = true
				fmt.Printf("Replacing existing %s installation...\n", service)
			}

			// Install the service (this will build to a new location)
			if err := installer.InstallWithReplace(wasInstalled); err != nil {
				utils.PrintError("Installation failed for %s: %v", service, err)
				failed = append(failed, service)
				continue
			}
		}

		fmt.Println()
		if len(failed) > 0 {
			fmt.Printf("❌ Some installations failed:\n")
			for _, service := range failed {
				fmt.Printf("   • %s\n", service)
			}
			fmt.Printf("\n💡 Run diagnostics: yerd doctor\n")
		} else {
			fmt.Printf("✓ All web services installed successfully\n")
		}
	},
}

func init() {
	InstallCmd.Flags().BoolP("force", "f", false, "Force reinstall even if services are already installed")
}
