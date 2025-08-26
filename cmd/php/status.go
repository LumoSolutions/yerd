package php

import (
	"fmt"

	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show YERD PHP status and configuration",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()

			all := make(config.PhpConfig)
			config.GetStruct("php", &all)

			outputYerdConfig(all)

			for _, version := range all {
				outputPhpInfo(version)
			}
		},
	}
}

func outputYerdConfig(fullConfig config.PhpConfig) {
	fmt.Printf("📊 YERD PHP Status\n")
	fmt.Printf("├─ Installed PHP versions: %d\n", len(fullConfig))

	cliVersion := ""
	for _, version := range fullConfig {
		if version.IsCLI {
			cliVersion = version.Version
		}
	}

	if cliVersion != "" {
		fmt.Printf("├─ Current CLI: PHP %s\n", cliVersion)
	} else {
		fmt.Printf("├─ Current CLI: None set\n")
	}

	fmt.Printf("└─ Config: ~/.config/yerd/config.json\n")
	fmt.Println()
}

func outputPhpInfo(info config.PhpInfo) {
	status := getServiceStatus(info.Version)
	flag := "🟢"
	if status != "Running" {
		flag = "🔴"
	}

	fmt.Printf("%s PHP %s (%s)\n", flag, info.Version, info.InstalledVersion)
	fmt.Printf("├─ Binary: %s\n", constants.YerdBinDir+fmt.Sprintf("/php%s", info.Version))
	fmt.Printf("├─ php.ini: %s\n", constants.YerdEtcDir+fmt.Sprintf("/php%s/php.ini", info.Version))
	fmt.Printf("├─ FPM Socket: %s\n", constants.YerdPHPDir+fmt.Sprintf("/run/php%s-fpm.sock", info.Version))
	fmt.Printf("└─ FPM Service: %s\n", status)
	fmt.Println()
}

func getServiceStatus(version string) string {
	if utils.SystemdServiceActive(fmt.Sprintf("yerd-php%s-fpm", version)) {
		return "Running"
	}
	return "Stopped"
}
