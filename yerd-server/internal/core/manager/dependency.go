package manager

import (
	"fmt"
	"strings"

	"github.com/lumosolutions/yerd/server/internal/constants"
	"github.com/lumosolutions/yerd/server/internal/core/services"
	"github.com/lumosolutions/yerd/server/internal/utils"
)

type DependencyManager struct {
	Distro    string
	Pm        string
	PmCommand string
	Cm        *CommandManager
	Log       *services.Logger
}

func NewDependencyManager(logger *services.Logger) (*DependencyManager, error) {
	cm := NewCommandManager(logger)

	distro, err := detectDistribution(cm)
	if err != nil {
		return nil, fmt.Errorf("failed to detect distribution: %v", err)
	}

	pm, pmCmd, err := detectPackageManager(cm)
	if err != nil {
		return nil, fmt.Errorf("failed to detect package manager: %v", err)
	}

	return &DependencyManager{
		Distro:    distro,
		Pm:        pm,
		PmCommand: pmCmd,
		Cm:        cm,
		Log:       logger,
	}, nil
}

func detectDistribution(cm *CommandManager) (string, error) {
	if output, success := cm.ExecuteCommand("cat", "/etc/os-release"); success {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				distro := strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
				return distro, nil
			}
		}
	}

	if output, success := cm.ExecuteCommand("lsb_release", "-si"); !success {
		distro := strings.ToLower(strings.TrimSpace(output))
		return distro, nil
	}

	releaseFiles := map[string]string{
		"/etc/redhat-release": "rhel",
		"/etc/debian_version": "debian",
		"/etc/arch-release":   "arch",
		"/etc/SuSE-release":   "opensuse",
		"/etc/alpine-release": "alpine",
	}

	for file, distro := range releaseFiles {
		if _, success := cm.ExecuteCommand("test", "-f", file); success {
			return distro, nil
		}
	}

	return "unknown", fmt.Errorf("could not detect distribution")
}

func (dm *DependencyManager) InstallPackages(packages []string) error {
	var toInstall []string
	for _, pkg := range packages {
		if pkgs, exists := constants.GetSystemPackages(pkg, dm.Pm); exists {
			toInstall = append(toInstall, pkgs...)
		}
	}

	if len(toInstall) == 0 {
		return nil
	}

	pmConfig, exists := constants.GetPackageManagerConfig(dm.Pm)
	if !exists {
		return fmt.Errorf("unsupported package manager: %s", dm.Pm)
	}

	if utils.IsRunningElevated() {
		args := append(pmConfig.InstallArgs, toInstall...)
		if output, success := dm.Cm.ExecuteCommand(dm.PmCommand, args...); !success {
			dm.Log.Info("deps", "%s", output)
			return fmt.Errorf("package installation failed")
		}

		return nil
	}

	var missing []string
	for _, pkg := range toInstall {
		if !dm.isPackageInstalled(pkg) {
			missing = append(missing, pkg)
		}
	}

	if len(missing) > 0 {
		dm.Log.Info("deps", "missing deps: %s", strings.Join(missing, ","))
		return fmt.Errorf("missing the following dependencies %s", strings.Join(missing, ","))
	}

	dm.Log.Info("deps", "no dependencies required to install")
	return nil

}

func (dm *DependencyManager) isPackageInstalled(pkg string) bool {
	config, exists := constants.GetPackageManagerConfig(dm.Pm)
	if !exists {
		return dm.checkLibrary(pkg)
	}

	args := append(config.QueryArgs, pkg)
	output, success := dm.Cm.ExecuteCommand(config.QueryCmd, args...)

	if !success {
		return false
	}

	if dm.Pm == constants.APT {
		return strings.Contains(output, "ii")
	}

	return true
}

func (dm *DependencyManager) checkLibrary(libName string) bool {
	paths := []string{"/usr/lib", "/usr/local/lib", "/opt/homebrew/lib", "/lib"}

	for _, path := range paths {
		if _, success := dm.Cm.ExecuteCommand("find", path, "-name", libName+"*", "-type", "f"); success {
			return true
		}
	}

	return false
}

func detectPackageManager(cm *CommandManager) (string, string, error) {
	for pm, config := range constants.PackageManagerConfigs {
		if _, exists := cm.CommandExists(config.CheckName); exists {
			return pm, config.Command, nil
		}
	}

	return "", "", fmt.Errorf("no supported package manager found")
}
