package composer

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/constants"
	internalComposer "github.com/lumosolutions/yerd/internal/installers/composer"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstalls the YERD managed Composer",
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()
			green := color.New(color.FgGreen)
			yellow := color.New(color.FgYellow)
			blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			if !utils.FileExists(constants.LocalComposerPath) {
				yellow.Printf("Composer is not installed\n")
				blue.Printf("- To install composer, please run:\n")
				blue.Printf("- 'sudo yerd composer install'\n\n")

				red.Printf("❌ Operation cancelled\n")
				return
			}

			if err := internalComposer.RemoveComposer(); err != nil {
				red.Printf("Composer failed to uninstall!\n")
				blue.Printf("- Error: %v\n\n", err)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			green.Printf("✓ Composer was uninstalled\n")
		},
	}
}
