package composer

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/constants"
	internalComposer "github.com/lumosolutions/yerd/internal/installers/composer"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Updates the YERD managed Composer to the latest version",
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

			if err := internalComposer.InstallComposer(); err != nil {
				red.Printf("Composer failed to update!\n")
				blue.Printf("- Error: %v\n\n", err)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			green.Printf("✓ Composer updated successfully\n")
			blue.Printf("- Type it out with: 'composer --version'\n")
		},
	}
}
