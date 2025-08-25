package sites

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Sets a configuration value for a given site",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			//green := color.New(color.FgGreen)
			//yellow := color.New(color.FgYellow)
			blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			if len(args) < 3 {
				red.Println("At least three arguments are required")
				blue.Println("- yerd sites set <name> <value> <site>")
				blue.Println("- Examples:")
				blue.Println("- 'sudo yerd sites set php 8.3 example.test'")
				return
			}

			setName := args[0]
			setValue := args[1]
			siteIdentifier := args[2]

			siteManager, err := manager.NewSiteManager()
			if err != nil {
				red.Println("Unable to create a site manager instance")
				red.Println("Are the web components installed?")
				blue.Println("- You can install the web components with:")
				blue.Println("- 'sudo yerd web install'")
				return
			}

			siteManager.SetValue(setName, setValue, siteIdentifier)
		},
	}
}
