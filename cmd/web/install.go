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
	Long: `Install nginx and dnsmasq for local development.

Examples:
  yerd web install       # Install if not already installed
  yerd web install -f    # Force reinstall even if already installed
  yerd web install -c    # Force download fresh config files from GitHub
  yerd web install -f -c # Force reinstall and download fresh configs`,
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintSplash()

		if !utils.CheckAndPromptForSudo("Web services installation", "install") {
			return
		}

		force, _ := cmd.Flags().GetBool("force")
		forceConfig, _ := cmd.Flags().GetBool("forceConfig")
		services := []string{"nginx", "dnsmasq"}

		if force {
			fmt.Printf("Force reinstalling web services:\n")
		} else {
			fmt.Printf("Installing web services:\n")
		}

		if forceConfig {
			fmt.Printf("Force downloading fresh configuration files from GitHub\n")
		}
		fmt.Printf("  • nginx 1.29.1   - High-performance HTTP server and reverse proxy\n")
		fmt.Printf("  • dnsmasq 2.91   - Lightweight DNS forwarder and DHCP server\n")
		fmt.Println()

		var failed []string

		for _, service := range services {
			if !force && !forceConfig && web.IsServiceInstalled(service) {
				utils.PrintWarning("Service '%s' is already installed, skipping", service)
				utils.PrintInfo("Use -f or --force to reinstall")
				continue
			}

			if !force && forceConfig && web.IsServiceInstalled(service) {
				fmt.Printf("Updating configuration for existing %s installation...\n", service)
			}

			installer, err := web.NewWebInstaller(service)
			if err != nil {
				utils.PrintError("Failed to create installer for %s: %v", service, err)
				failed = append(failed, service)
				continue
			}

			installer.SetForceConfig(forceConfig)

			if !force && forceConfig && web.IsServiceInstalled(service) {
				if err := installer.UpdateConfigOnly(); err != nil {
					utils.PrintError("Config update failed for %s: %v", service, err)
					failed = append(failed, service)
					continue
				}
			} else {
				var wasInstalled bool
				if force && web.IsServiceInstalled(service) {
					wasInstalled = true
					fmt.Printf("Replacing existing %s installation...\n", service)
				}

				if err := installer.InstallWithReplace(wasInstalled); err != nil {
					utils.PrintError("Installation failed for %s: %v", service, err)
					failed = append(failed, service)
					continue
				}
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
	InstallCmd.Flags().BoolP("forceConfig", "c", false, "Force download fresh configuration files from GitHub")
}
