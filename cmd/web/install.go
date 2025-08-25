package web

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/installers/nginx"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs any web components required for local development sites",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			//green := color.New(color.FgGreen)
			//yellow := color.New(color.FgYellow)
			//blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			installer, err := nginx.NewNginxInstaller(false, true)
			if err != nil {
				red.Printf("Install failed\n\n")
			}

			installer.Install()
		},
	}
}
