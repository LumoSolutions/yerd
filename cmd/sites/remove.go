package sites

import (
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

			if !utils.CheckAndPromptForSudo() {
				return
			}

			path := args[0]

			siteManager, _ := manager.NewSiteManager()
			siteManager.RemoveSite(path)
		},
	}
}
