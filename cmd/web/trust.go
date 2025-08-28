package web

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/config"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/manager"
	"github.com/lumosolutions/yerd/internal/utils"
	"github.com/lumosolutions/yerd/internal/version"
	"github.com/spf13/cobra"
)

func BuildTrustCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "trust",
		Short: "Attempts to refresh the YERD CA for Chrome",
		Run: func(cmd *cobra.Command, args []string) {
			version.PrintSplash()
			red := color.New(color.FgRed)
			green := color.New(color.FgGreen)

			if !utils.CheckAndPromptForSudo() {
				return
			}

			webConfig := config.GetWebConfig()

			if !webConfig.Installed {
				red.Println("YERD web components are not installed")
				return
			}

			cm := manager.NewCertificateManager()

			caPath := filepath.Join(constants.CertsDir, "ca")
			caFile := "yerd.crt"

			cm.ChromeUntrust()
			if err := cm.ChromeTrust(caPath, caFile); err != nil {
				red.Println("Unable to trust CA cert with chrome due to the following error:")
				red.Println(err)
				return
			}

			green.Println("Chrome Trust Updated")
		},
	}
}
