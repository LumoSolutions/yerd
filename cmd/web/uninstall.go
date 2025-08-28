package web

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/installers/nginx"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstalls the web components required for local development",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			green := color.New(color.FgGreen)
			red := color.New(color.FgRed)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			installer, err := nginx.NewNginxInstaller(false, true)
			if err != nil {
				red.Printf("Install failed\n\n")
			}

			installer.Uninstall()

			newConfig := &config.WebConfig{
				Installed: false,
			}

			config.SetStruct("web", newConfig)

			green.Println("Successfully uninstalled web components")
		},
	}
}
