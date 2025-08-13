package dependencies

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/LumoSolutions/yerd/internal/utils"
	"github.com/fatih/color"
)

type PackageManager string

const (
	APT    PackageManager = "apt"
	YUM    PackageManager = "yum"
	DNF    PackageManager = "dnf"
	PACMAN PackageManager = "pacman"
	ZYPPER PackageManager = "zypper"
	APKL   PackageManager = "apk"
)

type DependencyManager struct {
	distro    string
	pm        PackageManager
	pmCommand string
	quiet     bool
}

var extensionDependencies = map[string]map[PackageManager][]string{
	"curl": {
		APT:    []string{"libcurl4-openssl-dev"},
		YUM:    []string{"libcurl-devel"},
		DNF:    []string{"libcurl-devel"},
		PACMAN: []string{"curl"},
		ZYPPER: []string{"libcurl-devel"},
		APKL:   []string{"curl-dev"},
	},
	"openssl": {
		APT:    []string{"libssl-dev"},
		YUM:    []string{"openssl-devel"},
		DNF:    []string{"openssl-devel"},
		PACMAN: []string{"openssl"},
		ZYPPER: []string{"openssl-devel"},
		APKL:   []string{"openssl-dev"},
	},
	"zip": {
		APT:    []string{"libzip-dev"},
		YUM:    []string{"libzip-devel"},
		DNF:    []string{"libzip-devel"},
		PACMAN: []string{"libzip"},
		ZYPPER: []string{"libzip-devel"},
		APKL:   []string{"libzip-dev"},
	},
	"gd": {
		APT:    []string{"libgd-dev"},
		YUM:    []string{"gd-devel"},
		DNF:    []string{"gd-devel"},
		PACMAN: []string{"gd"},
		ZYPPER: []string{"gd-devel"},
		APKL:   []string{"gd-dev"},
	},
	"mysqli": {
		APT:    []string{"libmysqlclient-dev"},
		YUM:    []string{"mysql-devel"},
		DNF:    []string{"mysql-devel"},
		PACMAN: []string{"mariadb-libs"},
		ZYPPER: []string{"libmysqlclient-devel"},
		APKL:   []string{"mysql-dev"},
	},
	"pdo-mysql": {
		APT:    []string{"libmysqlclient-dev"},
		YUM:    []string{"mysql-devel"},
		DNF:    []string{"mysql-devel"},
		PACMAN: []string{"mariadb-libs"},
		ZYPPER: []string{"libmysqlclient-devel"},
		APKL:   []string{"mysql-dev"},
	},
	"pgsql": {
		APT:    []string{"libpq-dev"},
		YUM:    []string{"postgresql-devel"},
		DNF:    []string{"postgresql-devel"},
		PACMAN: []string{"postgresql-libs"},
		ZYPPER: []string{"postgresql-devel"},
		APKL:   []string{"postgresql-dev"},
	},
	"pdo-pgsql": {
		APT:    []string{"libpq-dev"},
		YUM:    []string{"postgresql-devel"},
		DNF:    []string{"postgresql-devel"},
		PACMAN: []string{"postgresql-libs"},
		ZYPPER: []string{"postgresql-devel"},
		APKL:   []string{"postgresql-dev"},
	},
	"jpeg": {
		APT:    []string{"libjpeg-dev"},
		YUM:    []string{"libjpeg-turbo-devel"},
		DNF:    []string{"libjpeg-turbo-devel"},
		PACMAN: []string{"libjpeg-turbo"},
		ZYPPER: []string{"libjpeg8-devel"},
		APKL:   []string{"libjpeg-turbo-dev"},
	},
	"freetype": {
		APT:    []string{"libfreetype6-dev"},
		YUM:    []string{"freetype-devel"},
		DNF:    []string{"freetype-devel"},
		PACMAN: []string{"freetype2"},
		ZYPPER: []string{"freetype2-devel"},
		APKL:   []string{"freetype-dev"},
	},
	"zlib": {
		APT:    []string{"zlib1g-dev"},
		YUM:    []string{"zlib-devel"},
		DNF:    []string{"zlib-devel"},
		PACMAN: []string{"zlib"},
		ZYPPER: []string{"zlib-devel"},
		APKL:   []string{"zlib-dev"},
	},
	"bz2": {
		APT:    []string{"libbz2-dev"},
		YUM:    []string{"bzip2-devel"},
		DNF:    []string{"bzip2-devel"},
		PACMAN: []string{"bzip2"},
		ZYPPER: []string{"libbz2-devel"},
		APKL:   []string{"bzip2-dev"},
	},
	"intl": {
		APT:    []string{"libicu-dev"},
		YUM:    []string{"libicu-devel"},
		DNF:    []string{"libicu-devel"},
		PACMAN: []string{"icu"},
		ZYPPER: []string{"libicu-devel"},
		APKL:   []string{"icu-dev"},
	},
	"gettext": {
		APT:    []string{"gettext"},
		YUM:    []string{"gettext-devel"},
		DNF:    []string{"gettext-devel"},
		PACMAN: []string{"gettext"},
		ZYPPER: []string{"gettext-tools"},
		APKL:   []string{"gettext-dev"},
	},
	"gmp": {
		APT:    []string{"libgmp-dev"},
		YUM:    []string{"gmp-devel"},
		DNF:    []string{"gmp-devel"},
		PACMAN: []string{"gmp"},
		ZYPPER: []string{"gmp-devel"},
		APKL:   []string{"gmp-dev"},
	},
	"ldap": {
		APT:    []string{"libldap2-dev"},
		YUM:    []string{"openldap-devel"},
		DNF:    []string{"openldap-devel"},
		PACMAN: []string{"libldap"},
		ZYPPER: []string{"openldap2-devel"},
		APKL:   []string{"openldap-dev"},
	},
	"pcre": {
		APT:    []string{"libpcre2-dev"},
		YUM:    []string{"pcre2-devel"},
		DNF:    []string{"pcre2-devel"},
		PACMAN: []string{"pcre2"},
		ZYPPER: []string{"pcre2-devel"},
		APKL:   []string{"pcre2-dev"},
	},
}

var buildDependencies = map[PackageManager][]string{
	APT:    []string{"build-essential", "autoconf", "pkg-config", "re2c", "libonig-dev", "libxml2-dev", "libsqlite3-dev"},
	YUM:    []string{"gcc", "gcc-c++", "make", "autoconf", "pkgconfig", "re2c", "oniguruma-devel", "libxml2-devel", "sqlite-devel"},
	DNF:    []string{"gcc", "gcc-c++", "make", "autoconf", "pkgconf", "re2c", "oniguruma-devel", "libxml2-devel", "sqlite-devel"},
	PACMAN: []string{"base-devel", "autoconf", "pkgconf", "re2c", "oniguruma", "libxml2", "sqlite"},
	ZYPPER: []string{"gcc", "gcc-c++", "make", "autoconf", "pkg-config", "re2c", "libonig-devel", "libxml2-devel", "sqlite3-devel"},
	APKL:   []string{"build-base", "autoconf", "pkgconf", "re2c", "oniguruma-dev", "libxml2-dev", "sqlite-dev"},
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

// detectDistribution identifies the Linux distribution using multiple detection methods.
// Returns distribution name or error if detection fails.
func detectDistribution() (string, error) {
	if output, err := utils.ExecuteCommand("cat", "/etc/os-release"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				return strings.Trim(strings.TrimPrefix(line, "ID="), `"`), nil
			}
		}
	}

	if output, err := utils.ExecuteCommand("lsb_release", "-si"); err == nil {
		return strings.ToLower(strings.TrimSpace(output)), nil
	}

	releaseFiles := map[string]string{
		"/etc/redhat-release": "rhel",
		"/etc/debian_version": "debian",
		"/etc/arch-release":   "arch",
		"/etc/SuSE-release":   "opensuse",
		"/etc/alpine-release": "alpine",
	}

	for file, distro := range releaseFiles {
		if _, err := utils.ExecuteCommand("test", "-f", file); err == nil {
			return distro, nil
		}
	}

	return "unknown", fmt.Errorf("could not detect distribution")
}

// detectPackageManager finds available package manager by checking system commands.
// Returns PackageManager type, command string, and error if none found.
func detectPackageManager() (PackageManager, string, error) {
	managers := []struct {
		pm      PackageManager
		command string
		check   string
	}{
		{APT, "apt-get", "apt-get"},
		{DNF, "dnf", "dnf"},
		{YUM, "yum", "yum"},
		{PACMAN, "pacman", "pacman"},
		{ZYPPER, "zypper", "zypper"},
		{APKL, "apk", "apk"},
	}

	for _, mgr := range managers {
		if _, err := exec.LookPath(mgr.check); err == nil {
			return mgr.pm, mgr.command, nil
		}
	}

	return "", "", fmt.Errorf("no supported package manager found")
}

// GetDistro returns the detected Linux distribution name.
func (dm *DependencyManager) GetDistro() string {
	return dm.distro
}

// GetPackageManager returns the detected package manager type.
func (dm *DependencyManager) GetPackageManager() PackageManager {
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
	deps, exists := buildDependencies[dm.pm]
	if !exists {
		return fmt.Errorf("no build dependencies defined for package manager: %s", dm.pm)
	}

	if !dm.quiet {
		color.New(color.FgYellow).Printf("Installing build dependencies for %s...\n", dm.distro)
	}

	return dm.installPackages(deps, "build dependencies")
}

// InstallExtensionDependencies installs libraries required for specific PHP extensions.
// extensions: List of PHP extensions needing dependencies. Returns error if installation fails.
func (dm *DependencyManager) InstallExtensionDependencies(extensions []string) error {
	var allPackages []string
	missingDeps := make(map[string]bool)

	for _, ext := range extensions {
		if deps, exists := extensionDependencies[ext]; exists {
			if packages, hasPM := deps[dm.pm]; hasPM {
				for _, pkg := range packages {
					if !contains(allPackages, pkg) {
						allPackages = append(allPackages, pkg)
					}
				}
			} else {
				missingDeps[ext] = true
			}
		}
	}

	if len(missingDeps) > 0 && !dm.quiet {
		var missing []string
		for dep := range missingDeps {
			missing = append(missing, dep)
		}
		color.New(color.FgYellow).Printf("Warning: No package mappings for dependencies: %s\n", strings.Join(missing, ", "))
	}

	if len(allPackages) == 0 {
		if !dm.quiet {
			color.New(color.FgGreen).Println("No additional dependencies required for selected extensions")
		}
		return nil
	}

	if !dm.quiet {
		// Removed verbose extension dependency listing for cleaner output
	}

	return dm.installPackages(allPackages, "extension dependencies")
}

// installPackages executes package installation commands for the detected package manager.
// packages: Package list to install, description: Operation description for logging. Returns error if installation fails.
func (dm *DependencyManager) installPackages(packages []string, description string) error {
	if len(packages) == 0 {
		return nil
	}

	var cmd *exec.Cmd

	switch dm.pm {
	case APT:
		args := append([]string{"install", "-y"}, packages...)
		cmd = exec.Command("apt-get", args...)
	case DNF:
		args := append([]string{"install", "-y"}, packages...)
		cmd = exec.Command("dnf", args...)
	case YUM:
		args := append([]string{"install", "-y"}, packages...)
		cmd = exec.Command("yum", args...)
	case PACMAN:
		args := append([]string{"-S", "--noconfirm"}, packages...)
		cmd = exec.Command("pacman", args...)
	case ZYPPER:
		args := append([]string{"install", "-y"}, packages...)
		cmd = exec.Command("zypper", args...)
	case APKL:
		args := append([]string{"add"}, packages...)
		cmd = exec.Command("apk", args...)
	default:
		return fmt.Errorf("unsupported package manager: %s", dm.pm)
	}

	if !dm.quiet {
		// Removed verbose command output for cleaner installation experience
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.New(color.FgRed).Printf("Failed to install %s:\n%s\n", description, string(output))
		return fmt.Errorf("package installation failed: %v", err)
	}

	if !dm.quiet {
		color.New(color.FgGreen).Printf("✓ Successfully installed %s\n", description)
	}
	return nil
}


// contains checks if a string slice contains a specific item.
// slice: String slice to search, item: Item to find. Returns true if found.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// CheckSystemDependencies verifies which extension dependencies are missing from the system.
// extensions: List of PHP extensions to check. Returns slice of missing dependency names.
func (dm *DependencyManager) CheckSystemDependencies(extensions []string) []string {
	var missing []string

	for _, ext := range extensions {
		if _, exists := extensionDependencies[ext]; exists {
			if !dm.isDependencyAvailable(ext) {
				missing = append(missing, ext)
			}
		}
	}

	return missing
}

// isDependencyAvailable checks if a system dependency is available using pkg-config and custom checks.
// depName: Dependency name to check. Returns true if available.
func (dm *DependencyManager) isDependencyAvailable(depName string) bool {
	pkgConfigNames := map[string][]string{
		"gd":        {"gdlib"},
		"zip":       {"libzip"},
		"curl":      {"libcurl"},
		"openssl":   {"openssl"},
		"zlib":      {"zlib"},
		"libxml2":   {"libxml-2.0"},
		"freetype2": {"freetype2"},
		"icu":       {"icu-uc", "icu-io"},
		"pcre2":     {"libpcre2-8"},
		"bzip2":     {"bzip2"},
		"openldap":  {"ldap"},
	}

	if pkgNames, exists := pkgConfigNames[depName]; exists {
		for _, pkgName := range pkgNames {
			if _, err := exec.Command("pkg-config", "--exists", pkgName).CombinedOutput(); err == nil {
				return true
			}
		}
	}

	if _, err := exec.Command("pkg-config", "--exists", depName).CombinedOutput(); err == nil {
		return true
	}

	specialChecks := map[string]func() bool{
		"mysql": func() bool {
			return dm.checkCommand("mysql_config") || dm.checkLibrary("libmysqlclient")
		},
		"postgresql": func() bool {
			return dm.checkCommand("pg_config") || dm.checkLibrary("libpq")
		},
		"sqlite": func() bool {
			return dm.checkLibrary("libsqlite3") || dm.checkCommand("sqlite3")
		},
		"gmp": func() bool {
			return dm.checkLibrary("libgmp")
		},
		"gettext": func() bool {
			return dm.checkCommand("gettext") || dm.checkLibrary("libintl")
		},
	}

	if checker, exists := specialChecks[depName]; exists {
		return checker()
	}

	return false
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
		if _, err := utils.ExecuteCommand("find", path, "-name", libName+"*", "-type", "f"); err == nil {
			return true
		}
	}

	return false
}
