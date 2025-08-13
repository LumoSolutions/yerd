package installer

import (
	"fmt"

	"github.com/LumoSolutions/yerd/internal/dependencies"
	"github.com/LumoSolutions/yerd/internal/utils"
)

func checkBuildDependencies(logger *utils.Logger) error {
	utils.SafeLog(logger, "Checking and installing build dependencies...")

	depMgr, err := dependencies.NewDependencyManager()
	if err != nil {
		utils.SafeLog(logger, "Failed to initialize dependency manager: %v", err)
		return fmt.Errorf("failed to initialize dependency manager: %v", err)
	}

	utils.SafeLog(logger, "Detected %s with %s package manager", depMgr.GetDistro(), depMgr.GetPackageManager())
	fmt.Printf("📋 Detected: %s with %s\n", depMgr.GetDistro(), depMgr.GetPackageManager())

	if err := depMgr.InstallBuildDependencies(); err != nil {
		utils.SafeLog(logger, "Failed to install build dependencies: %v", err)
		return fmt.Errorf("failed to install build dependencies: %v", err)
	}

	utils.SafeLog(logger, "Build dependencies installed successfully")
	return nil
}
