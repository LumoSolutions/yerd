package sites

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Removes a local development site given a directory",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			red := color.New(color.FgRed)
			blue := color.New(color.FgBlue)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			path := args[0]

			siteManager, err := manager.NewSiteManager()
			if err != nil {
				red.Println("Unable to create a site manager instance")
				red.Println("Are the web components installed?")
				blue.Println("- You can install the web components with:")
				blue.Println("- 'sudo yerd web install'")
				return
			}

			siteManager.RemoveSite(path)
		},
	}
}
