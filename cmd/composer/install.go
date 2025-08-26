package composer

import (
	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/constants"
	internalComposer "github.com/lumosolutions/yerd/internal/installers/composer"
	"github.com/lumosolutions/yerd/internal/utils"
	intVersion "github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs a YERD managed Composer",
		Run: func(cmd *cobra.Command, args []string) {
			intVersion.PrintSplash()
			green := color.New(color.FgGreen)
			yellow := color.New(color.FgYellow)
			blue := color.New(color.FgBlue)
			red := color.New(color.FgRed)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			if utils.FileExists(constants.LocalComposerPath) {
				yellow.Printf("Composer is already installed\n")
				blue.Printf("- To upgrade composer, please run:\n")
				blue.Printf("- 'sudo yerd composer update'\n\n")

				red.Printf("❌ Operation cancelled\n")
				return
			}

			if err := internalComposer.InstallComposer(); err != nil {
				red.Printf("Composer failed to install!\n")
				blue.Printf("- Error: %v\n\n", err)
				red.Printf("❌ Operation cancelled\n")
				return
			}

			green.Printf("✓ Composer installed successfully\n")
			blue.Printf("- Type it out with: 'composer --version'\n")
		},
	}
}
