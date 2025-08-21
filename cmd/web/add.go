package web

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Adds a new local development site given a directory",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			green := color.New(color.FgGreen)
			//yellow := color.New(color.FgYellow)
			//blue := color.New(color.FgBlue)
			//red := color.New(color.FgRed)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			green.Printf("✓ Web Components Installed Successfully\n")
		},
	}
}
