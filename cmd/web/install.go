package web

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/web"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install [service]",
	Short: "Install web services",
	Long:  `Install web services like nginx, apache, and other web servers.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		if !web.IsValidService(service) {
			fmt.Printf("Error: '%s' is not a supported web service\n", service)
			fmt.Println("Supported services:", web.GetSupportedServices())
			return
		}

		manager, err := web.NewWebManager()
		if err != nil {
			fmt.Printf("Error initializing web manager: %v\n", err)
			return
		}

		err = manager.InstallService(service)
		if err != nil {
			fmt.Printf("Error installing %s: %v\n", service, err)
			return
		}

		fmt.Printf("Successfully installed %s\n", service)
	},
}
