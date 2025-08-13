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
	"mysql": {
		APT:    []string{"libmysqlclient-dev"},
		YUM:    []string{"mysql-devel"},
		DNF:    []string{"mysql-devel"},
		PACMAN: []string{"mariadb-libs"},
		ZYPPER: []string{"libmysqlclient-devel"},
		APKL:   []string{"mysql-dev"},
	},
	"postgresql": {
		APT:    []string{"libpq-dev"},
		YUM:    []string{"postgresql-devel"},
		DNF:    []string{"postgresql-devel"},
		PACMAN: []string{"postgresql-libs"},
		ZYPPER: []string{"postgresql-devel"},
		APKL:   []string{"postgresql-dev"},
	},
	"sqlite": {
		APT:    []string{"libsqlite3-dev"},
		YUM:    []string{"sqlite-devel"},
		DNF:    []string{"sqlite-devel"},
		PACMAN: []string{"sqlite"},
		ZYPPER: []string{"sqlite3-devel"},
		APKL:   []string{"sqlite-dev"},
	},
	"libjpeg": {
		APT:    []string{"libjpeg-dev"},
		YUM:    []string{"libjpeg-turbo-devel"},
		DNF:    []string{"libjpeg-turbo-devel"},
		PACMAN: []string{"libjpeg-turbo"},
		ZYPPER: []string{"libjpeg8-devel"},
		APKL:   []string{"libjpeg-turbo-dev"},
	},
	"freetype2": {
		APT:    []string{"libfreetype6-dev"},
		YUM:    []string{"freetype-devel"},
		DNF:    []string{"freetype-devel"},
		PACMAN: []string{"freetype2"},
		ZYPPER: []string{"freetype2-devel"},
		APKL:   []string{"freetype-dev"},
	},
	"libxml2": {
		APT:    []string{"libxml2-dev"},
		YUM:    []string{"libxml2-devel"},
		DNF:    []string{"libxml2-devel"},
		PACMAN: []string{"libxml2"},
		ZYPPER: []string{"libxml2-devel"},
		APKL:   []string{"libxml2-dev"},
	},
	"zlib": {
		APT:    []string{"zlib1g-dev"},
		YUM:    []string{"zlib-devel"},
		DNF:    []string{"zlib-devel"},
		PACMAN: []string{"zlib"},
		ZYPPER: []string{"zlib-devel"},
		APKL:   []string{"zlib-dev"},
	},
	"bzip2": {
		APT:    []string{"libbz2-dev"},
		YUM:    []string{"bzip2-devel"},
		DNF:    []string{"bzip2-devel"},
		PACMAN: []string{"bzip2"},
		ZYPPER: []string{"libbz2-devel"},
		APKL:   []string{"bzip2-dev"},
	},
	"icu": {
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
	"openldap": {
		APT:    []string{"libldap2-dev"},
		YUM:    []string{"openldap-devel"},
		DNF:    []string{"openldap-devel"},
		PACMAN: []string{"libldap"},
		ZYPPER: []string{"openldap2-devel"},
		APKL:   []string{"openldap-dev"},
	},
	"pcre2": {
		APT:    []string{"libpcre2-dev"},
		YUM:    []string{"pcre2-devel"},
		DNF:    []string{"pcre2-devel"},
		PACMAN: []string{"pcre2"},
		ZYPPER: []string{"pcre2-devel"},
		APKL:   []string{"pcre2-dev"},
	},
}

var buildDependencies = map[PackageManager][]string{
	APT:    []string{"build-essential", "autoconf", "pkg-config", "re2c", "libonig-dev"},
	YUM:    []string{"gcc", "gcc-c++", "make", "autoconf", "pkgconfig", "re2c", "oniguruma-devel"},
	DNF:    []string{"gcc", "gcc-c++", "make", "autoconf", "pkgconf", "re2c", "oniguruma-devel"},
	PACMAN: []string{"base-devel", "autoconf", "pkgconf", "re2c", "oniguruma"},
	ZYPPER: []string{"gcc", "gcc-c++", "make", "autoconf", "pkg-config", "re2c", "libonig-devel"},
	APKL:   []string{"build-base", "autoconf", "pkgconf", "re2c", "oniguruma-dev"},
}

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
	}, nil
}

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

func (dm *DependencyManager) GetDistro() string {
	return dm.distro
}

func (dm *DependencyManager) GetPackageManager() PackageManager {
	return dm.pm
}

func (dm *DependencyManager) InstallBuildDependencies() error {
	deps, exists := buildDependencies[dm.pm]
	if !exists {
		return fmt.Errorf("no build dependencies defined for package manager: %s", dm.pm)
	}

	color.New(color.FgYellow).Printf("Installing build dependencies for %s...\n", dm.distro)

	return dm.installPackages(deps, "build dependencies")
}

func (dm *DependencyManager) InstallExtensionDependencies(extensions []string) error {
	var allPackages []string
	missingDeps := make(map[string]bool)

	for _, ext := range extensions {
		depName := mapExtensionToDependency(ext)
		if depName == "" {
			continue
		}

		if deps, exists := extensionDependencies[depName]; exists {
			if packages, hasPM := deps[dm.pm]; hasPM {
				for _, pkg := range packages {
					if !contains(allPackages, pkg) {
						allPackages = append(allPackages, pkg)
					}
				}
			} else {
				missingDeps[depName] = true
			}
		}
	}

	if len(missingDeps) > 0 {
		var missing []string
		for dep := range missingDeps {
			missing = append(missing, dep)
		}
		color.New(color.FgYellow).Printf("Warning: No package mappings for dependencies: %s\n", strings.Join(missing, ", "))
	}

	if len(allPackages) == 0 {
		color.New(color.FgGreen).Println("No additional dependencies required for selected extensions")
		return nil
	}

	color.New(color.FgCyan).Printf("Installing dependencies for extensions: %s\n", strings.Join(extensions, ", "))

	return dm.installPackages(allPackages, "extension dependencies")
}

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

	color.New(color.FgBlue).Printf("Running: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.New(color.FgRed).Printf("Failed to install %s:\n%s\n", description, string(output))
		return fmt.Errorf("package installation failed: %v", err)
	}

	color.New(color.FgGreen).Printf("✓ Successfully installed %s\n", description)
	return nil
}

func mapExtensionToDependency(extension string) string {
	mapping := map[string]string{
		"curl":       "curl",
		"openssl":    "openssl",
		"zip":        "zip",
		"gd":         "gd",
		"mysqli":     "mysql",
		"pdo-mysql":  "mysql",
		"pgsql":      "postgresql",
		"pdo-pgsql":  "postgresql",
		"sqlite3":    "sqlite",
		"pdo-sqlite": "sqlite",
		"jpeg":       "libjpeg",
		"freetype":   "freetype2",
		"xml":        "libxml2",
		"zlib":       "zlib",
		"bz2":        "bzip2",
		"intl":       "icu",
		"gettext":    "gettext",
		"gmp":        "gmp",
		"ldap":       "openldap",
		"pcre":       "pcre2",
	}

	return mapping[extension]
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (dm *DependencyManager) CheckSystemDependencies(extensions []string) []string {
	var missing []string

	for _, ext := range extensions {
		depName := mapExtensionToDependency(ext)
		if depName == "" {
			continue
		}

		if !dm.isDependencyAvailable(depName) {
			missing = append(missing, depName)
		}
	}

	return missing
}

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

func (dm *DependencyManager) checkCommand(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func (dm *DependencyManager) checkLibrary(libName string) bool {
	paths := []string{"/usr/lib", "/usr/local/lib", "/opt/homebrew/lib", "/lib"}

	for _, path := range paths {
		if _, err := utils.ExecuteCommand("find", path, "-name", libName+"*", "-type", "f"); err == nil {
			return true
		}
	}

	return false
}
