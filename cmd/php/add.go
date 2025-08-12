package php

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/LumoSolutions/yerd/internal/installer"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/pkg/php"
)

var AddCmd = &cobra.Command{
	Use:   "add [version]",
	Short: "Install a specific PHP version",
	Long: `Build and install a PHP version from official source code.
	
Examples:
  yerd php add 8.3
  yerd php add 8.4
  yerd php add php8.3
  yerd php add 8.4 -u    # bypass cache for latest version info`,
	Args: cobra.ExactArgs(1),
	Run:  runAdd,
}

func init() {
	AddCmd.Flags().BoolP("uncached", "u", false, "Bypass cache to get the latest version information from PHP.net")
}

func runAdd(cmd *cobra.Command, args []string) {
	version.PrintSplash()
	
	if !utils.CheckAndPromptForSudo("PHP installation", "add", args[0]) {
		return
	}
	
	versionArg := args[0]
	phpVersion := utils.NormalizePHPVersion(versionArg)
	
	if !php.IsValidVersion(phpVersion) {
		fmt.Printf("Error: Invalid PHP version '%s'\n", phpVersion)
		fmt.Printf("Available versions: %s\n", strings.Join(php.GetAvailableVersions(), ", "))
		return
	}
	
	uncached, _ := cmd.Flags().GetBool("uncached")
	
	fmt.Printf("Installing PHP %s...\n", phpVersion)
	if uncached {
		fmt.Printf("ℹ️  Bypassing cache to get latest version information\n")
	}
	
	err := installer.InstallPHP(phpVersion, uncached)
	if err != nil {
		fmt.Printf("\n❌ Installation failed: %v\n", err)
		
		fmt.Printf("💡 Run diagnostics: yerd doctor php%s\n", phpVersion)
		return
	}
	
	fmt.Printf("✓ PHP %s installed successfully\n", phpVersion)
}