package sites

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists development sites & their configuration",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			red := color.New(color.FgRed)

			sm, err := manager.NewSiteManager()
			if err != nil {
				red.Printf("Unable to create an instance of SiteManager\n")
				return
			}

			sm.ListSites()
		},
	}
}
