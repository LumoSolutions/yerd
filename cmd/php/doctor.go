package php

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/LumoSolutions/yerd/internal/version"
	"github.com/LumoSolutions/yerd/pkg/php"
)

var DoctorCmd = &cobra.Command{
	Use:   "doctor [version]",
	Short: "Diagnose PHP installation issues",
	Long:  `Run diagnostics to help troubleshoot PHP installation and configuration issues.

Examples:
  yerd php doctor
  yerd php doctor 8.3
  yerd php doctor php8.4`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) {
	version.PrintSplash()
	
	fmt.Printf("🩺 YERD Doctor - System Diagnostics\n\n")
	
	fmt.Printf("1️⃣  System Requirements\n")
	sysReq := utils.CheckSystemRequirements()
	utils.PrintSystemRequirements(sysReq)
	fmt.Println()
	
	fmt.Printf("2️⃣  YERD Configuration\n")
	yerdConfig := utils.CheckYERDConfiguration()
	utils.PrintYERDConfiguration(yerdConfig)
	fmt.Println()
	
	fmt.Printf("3️⃣  System PHP Conflicts\n")
	phpConflicts := utils.CheckSystemPHPConflicts()
	utils.PrintSystemPHPConflicts(phpConflicts)
	fmt.Println()
	
	if len(args) > 0 {
		phpVersion := utils.NormalizePHPVersion(args[0])
		fmt.Printf("4️⃣  PHP %s Diagnostics\n", phpVersion)
		availableVersions := php.GetAvailableVersions()
		phpDiag := utils.DiagnosePHPVersion(phpVersion, availableVersions)
		utils.PrintPHPVersionDiagnostics(phpDiag, availableVersions)
		fmt.Println()
		
		fmt.Printf("5️⃣  Installed PHP Binaries\n")
		binaries := utils.FindInstalledPHPBinaries()
		utils.PrintInstalledPHPBinaries(binaries)
		fmt.Println()
	} else {
		fmt.Printf("4️⃣  Installed PHP Binaries\n")
		binaries := utils.FindInstalledPHPBinaries()
		utils.PrintInstalledPHPBinaries(binaries)
		fmt.Println()
	}
	
	fmt.Printf("✅ Diagnostics complete. Use this information to troubleshoot issues.\n")
}





