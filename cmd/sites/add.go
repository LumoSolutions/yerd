package sites

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds a new local development site given a directory",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			red := color.New(color.FgRed)
			blue := color.New(color.FgBlue)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			path := args[0]
			domain, _ := cmd.Flags().GetString("domain")
			folder, _ := cmd.Flags().GetString("folder")
			php, _ := cmd.Flags().GetString("php")

			siteManager, err := manager.NewSiteManager()
			if err != nil {
				red.Println("Unable to create a site manager instance")
				red.Println("Are the web components installed?")
				blue.Println("- You can install the web components with:")
				blue.Println("- 'sudo yerd web install'")
				return
			}

			siteManager.AddSite(path, domain, folder, php)
		},
	}

	cmd.Flags().StringP("domain", "d", "", "Override the default domain value (eg: mysite.test)")
	cmd.Flags().StringP("folder", "f", "", "Specify a public directory under the root")
	cmd.Flags().StringP("php", "p", "", "Specify the version of php to use")

	return cmd
}
