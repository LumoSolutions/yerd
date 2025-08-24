package manager

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

const (
	APT    = constants.APT
	YUM    = constants.YUM
	DNF    = constants.DNF
	PACMAN = constants.PACMAN
	ZYPPER = constants.ZYPPER
	APKL   = constants.APKL
)

type DependencyManager struct {
	distro    string
	pm        string
	pmCommand string
	quiet     bool
}

// NewDependencyManager creates a new dependency manager with auto-detected distribution and package manager.
// Returns configured DependencyManager or error if detection fails.
func NewDependencyManager() (*DependencyManager, error) {
	distro, err := detectDistribution()
	if err != nil {
		return nil, fmt.Errorf("failed to detect distribution: %v", err)
	}

	pm, pmCmd, err := detectPackageManager()
	if err != nil {
		return nil, fmt.Errorf("failed to detect package manager: %v", err)
	}

	return &DependencyManager{
		distro:    distro,
		pm:        pm,
		pmCommand: pmCmd,
		quiet:     false,
	}, nil
}

func (dm *DependencyManager) TrustCertificate(certificate, name string) error {
	switch dm.distro {
	case "ubuntu", "debian":
		if _, success := utils.ExecuteCommand("cp", certificate,
			fmt.Sprintf("/usr/local/share/ca-certificates/%s-ca.crt", name)); !success {
			return fmt.Errorf("failed to copy certificate for %s", dm.distro)
		}
		if _, success := utils.ExecuteCommand("update-ca-certificates"); !success {
			return fmt.Errorf("failed to update ca-certificates for %s", dm.distro)
		}
		break

	case "arch", "manjaro":
		if _, success := utils.ExecuteCommand("cp", certificate,
			fmt.Sprintf("/etc/ca-certificates/trust-source/anchors/%s-ca.crt", name)); !success {
			return fmt.Errorf("failed to copy certificate for %s", dm.distro)
		}
		if _, success := utils.ExecuteCommand("trust", "extract-compat"); !success {
			if _, success := utils.ExecuteCommand("update-ca-trust"); !success {
				return fmt.Errorf("failed to update trust store for %s", dm.distro)
			}
		}
		break

	case "rhel", "centos", "fedora", "rocky", "almalinux":
		if _, success := utils.ExecuteCommand("cp", certificate,
			fmt.Sprintf("/etc/pki/ca-trust/source/anchors/%s-ca.crt", name)); !success {
			return fmt.Errorf("failed to copy certificate for %s", dm.distro)
		}
		if _, success := utils.ExecuteCommand("update-ca-trust"); !success {
			return fmt.Errorf("failed to update ca-trust for %s", dm.distro)
		}
		break

	case "opensuse", "opensuse-leap", "opensuse-tumbleweed", "sles":
		if _, success := utils.ExecuteCommand("cp", certificate,
			fmt.Sprintf("/etc/pki/trust/anchors/%s-ca.crt", name)); !success {
			return fmt.Errorf("failed to copy certificate for %s", dm.distro)
		}
		if _, success := utils.ExecuteCommand("update-ca-certificates"); !success {
			return fmt.Errorf("failed to update ca-certificates for %s", dm.distro)
		}
		break

	case "alpine":
		if _, success := utils.ExecuteCommand("cp", certificate,
			fmt.Sprintf("/usr/local/share/ca-certificates/%s-ca.crt", name)); !success {
			return fmt.Errorf("failed to copy certificate for %s", dm.distro)
		}
		if _, success := utils.ExecuteCommand("update-ca-certificates"); !success {
			return fmt.Errorf("failed to update ca-certificates for %s", dm.distro)
		}
		break

	default:
		return fmt.Errorf("distro '%s' not supported yet", dm.distro)
	}

	return nil
}

// detectDistribution identifies the Linux distribution using multiple detection methods.
// Returns distribution name or error if detection fails.
func detectDistribution() (string, error) {
	if output, success := utils.ExecuteCommand("cat", "/etc/os-release"); success {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				distro := strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
				return distro, nil
			}
		}
	}

	if output, success := utils.ExecuteCommand("lsb_release", "-si"); !success {
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
		if _, success := utils.ExecuteCommand("test", "-f", file); success {
			return distro, nil
		}
	}

	return "unknown", fmt.Errorf("could not detect distribution")
}

// detectPackageManager finds available package manager by checking system commands.
// Returns PackageManager type, command string, and error if none found.
func detectPackageManager() (string, string, error) {
	for pm, config := range constants.PackageManagerConfigs {
		if _, exists := utils.CommandExists(config.CheckName); exists {
			return pm, config.Command, nil
		}
	}

	return "", "", fmt.Errorf("no supported package manager found")
}

// GetDistro returns the detected Linux distribution name.
func (dm *DependencyManager) GetDistro() string {
	return dm.distro
}

// GetPackageManager returns the detected package manager type.
func (dm *DependencyManager) GetPackageManager() string {
	return dm.pm
}

// SetQuiet enables or disables quiet mode for dependency installation.
// In quiet mode, verbose output is suppressed during operations.
func (dm *DependencyManager) SetQuiet(quiet bool) {
	dm.quiet = quiet
}

// InstallBuildDependencies installs essential build tools and dependencies for PHP compilation.
// Returns error if installation fails.
func (dm *DependencyManager) InstallBuildDependencies() error {
	deps := constants.GetBuildDependencies(dm.pm)
	if len(deps) == 0 {
		return fmt.Errorf("no build dependencies defined for package manager: %s", dm.pm)
	}

	return dm.installPackages(deps, "build dependencies")
}

// InstallWebBuildDependencies installs essential build tools and dependencies for web service compilation.
// Returns error if installation fails.
func (dm *DependencyManager) InstallWebDependencies() error {
	deps := constants.GetWebBuildDependencies(dm.pm)

	if len(deps) == 0 {
		return fmt.Errorf("no web build dependencies defined for package manager: %s", dm.pm)
	}

	return dm.installPackages(deps, "web build dependencies")
}

// InstallExtensionDependencies installs libraries required for specific PHP extensions.
// extensions: List of PHP extensions needing dependencies. Returns error if installation fails.
func (dm *DependencyManager) InstallExtensionDependencies(extensions []string) error {
	packages := dm.collectUniquePackages(extensions)

	if len(packages) == 0 {
		dm.logNoPackagesNeeded()
		return nil
	}

	return dm.installPackages(packages, "extension dependencies")
}

// collectUniquePackages gathers unique system packages for the given extensions
func (dm *DependencyManager) collectUniquePackages(extensions []string) []string {
	packageSet := make(map[string]bool)

	for _, ext := range extensions {
		if systemPkgs, exists := constants.GetSystemPackages(ext, dm.pm); exists {
			for _, pkg := range systemPkgs {
				packageSet[pkg] = true
			}
		}
	}

	return dm.mapKeysToSlice(packageSet)
}

// mapKeysToSlice converts map keys to slice
func (dm *DependencyManager) mapKeysToSlice(packageSet map[string]bool) []string {
	packages := make([]string, 0, len(packageSet))
	for pkg := range packageSet {
		packages = append(packages, pkg)
	}
	return packages
}

// logNoPackagesNeeded logs when no packages are required
func (dm *DependencyManager) logNoPackagesNeeded() {
	if !dm.quiet {
		color.New(color.FgGreen).Println("No additional dependencies required for selected extensions")
	}
}

// installPackages executes package installation commands for the detected package manager.
// packages: Package list to install, description: Operation description for logging. Returns error if installation fails.
func (dm *DependencyManager) installPackages(packages []string, description string) error {
	if len(packages) == 0 {
		return nil
	}

	config, exists := constants.GetPackageManagerConfig(dm.pm)
	if !exists {
		return fmt.Errorf("unsupported package manager: %s", dm.pm)
	}

	args := append(config.InstallArgs, packages...)
	cmd := exec.Command(dm.pmCommand, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogInfo("deps", "%s", string(output))
		return fmt.Errorf("package installation failed: %v", err)
	}

	return nil
}

// CheckSystemDependencies verifies which extension dependencies are missing from the system.
// extensions: List of PHP extensions to check. Returns slice of missing dependency names.
func (dm *DependencyManager) CheckSystemDependencies(extensions []string) []string {
	var missing []string

	for _, ext := range extensions {
		if _, exists := constants.GetDependencyConfig(ext); exists {
			if !dm.isDependencyAvailable(ext) {
				missing = append(missing, ext)
			}
		}
	}

	return missing
}

// isDependencyAvailable checks if a system dependency is available using distro-specific package detection and pkg-config.
// depName: Dependency name to check. Returns true if available.
func (dm *DependencyManager) isDependencyAvailable(depName string) bool {
	// Check centralized config first
	if dm.checkConfiguredDependency(depName) {
		return true
	}

	// Fallback to legacy pkg-config checks
	return dm.checkPkgConfig(depName)
}

// checkConfiguredDependency checks dependencies defined in the centralized config
func (dm *DependencyManager) checkConfiguredDependency(depName string) bool {
	config, exists := constants.GetDependencyConfig(depName)
	if !exists {
		return false
	}

	return dm.checkSystemPackages(config) ||
		dm.checkCommands(config) ||
		dm.checkLibraries(config) ||
		dm.checkCommonPkgConfig(config)
}

// checkSystemPackages checks if distro-specific packages are installed
func (dm *DependencyManager) checkSystemPackages(config *constants.DependencyConfig) bool {
	distroPackages, hasDistro := config.SystemPackages[dm.pm]
	if !hasDistro {
		return false
	}

	for _, pkgName := range distroPackages {
		if dm.isPackageInstalled(pkgName) {
			return true
		}
	}
	return false
}

// checkCommands checks if required commands are available
func (dm *DependencyManager) checkCommands(config *constants.DependencyConfig) bool {
	for _, cmd := range config.Commands {
		if dm.checkCommand(cmd) {
			return true
		}
	}
	return false
}

// checkLibraries checks if required libraries are available
func (dm *DependencyManager) checkLibraries(config *constants.DependencyConfig) bool {
	for _, lib := range config.Libraries {
		if dm.checkLibrary(lib) {
			return true
		}
	}
	return false
}

// checkCommonPkgConfig checks common pkg-config names
func (dm *DependencyManager) checkCommonPkgConfig(config *constants.DependencyConfig) bool {
	for _, pkgName := range config.CommonPkgConfig {
		if dm.checkPkgConfigName(pkgName) {
			return true
		}
	}
	return false
}

// checkPkgConfig performs fallback pkg-config checks
func (dm *DependencyManager) checkPkgConfig(depName string) bool {
	// Check legacy pkg-config names
	pkgConfigNames := dm.getPkgConfigNames()
	if pkgNames, exists := pkgConfigNames[depName]; exists {
		for _, pkgName := range pkgNames {
			if dm.checkPkgConfigName(pkgName) {
				return true
			}
		}
	}

	return dm.checkPkgConfigName(depName)
}

// checkPkgConfigName checks a single pkg-config name
func (dm *DependencyManager) checkPkgConfigName(pkgName string) bool {
	_, err := exec.Command("pkg-config", "--exists", pkgName).CombinedOutput()
	return err == nil
}

// checkCommand verifies if a system command is available in PATH.
// command: Command name to check. Returns true if command exists.
func (dm *DependencyManager) checkCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// checkLibrary searches common system library paths for a specific library file.
// libName: Library name to find. Returns true if library found in system paths.
func (dm *DependencyManager) checkLibrary(libName string) bool {
	paths := []string{"/usr/lib", "/usr/local/lib", "/opt/homebrew/lib", "/lib"}

	for _, path := range paths {
		if _, success := utils.ExecuteCommand("find", path, "-name", libName+"*", "-type", "f"); success {
			return true
		}
	}

	return false
}

// getPkgConfigNames returns distro-specific pkg-config package names for dependency detection.
// Returns map where keys are extension names and values are pkg-config package names to check.
func (dm *DependencyManager) getPkgConfigNames() map[string][]string {
	result := make(map[string][]string)

	// Get all dependencies from centralized config
	for _, depName := range constants.GetAllDependencyNames() {
		if config, exists := constants.GetDependencyConfig(depName); exists {
			if len(config.CommonPkgConfig) > 0 {
				result[depName] = append(result[depName], config.CommonPkgConfig...)
			}

			if distroNames, hasDistro := config.PkgConfigNames[dm.pm]; hasDistro {
				result[depName] = append(result[depName], distroNames...)
			}
		}
	}

	return result
}

// isPackageInstalled checks if a specific package is installed using the appropriate package manager.
// pkgName: Package name to check. Returns true if package is installed.
func (dm *DependencyManager) isPackageInstalled(pkgName string) bool {
	config, exists := constants.GetPackageManagerConfig(dm.pm)
	if !exists {
		return dm.checkLibrary(pkgName)
	}

	args := append(config.QueryArgs, pkgName)
	output, success := utils.ExecuteCommand(config.QueryCmd, args...)

	if !success {
		return false
	}

	if dm.pm == APT {
		return strings.Contains(output, "ii")
	}

	return true
}
