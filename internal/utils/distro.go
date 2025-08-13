package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const aptGetCmd = "apt-get"

type DistroInfo struct {
	Name           string
	PackageManager string
	InstallCmd     []string
	QueryCmd       []string
}

type PackageMapping struct {
	Arch   string
	Debian string
	RedHat string
	SUSE   string
}

// DetectDistro identifies the Linux distribution and package manager configuration.
// Returns DistroInfo with package management details or error if unsupported distribution.
func DetectDistro() (*DistroInfo, error) {
	if distro := detectFromOSRelease(); distro != nil {
		return distro, nil
	}

	return detectFromPackageManagers()
}

// detectFromOSRelease attempts to identify distribution from /etc/os-release file.
// Returns DistroInfo if identification succeeds, nil otherwise.
func detectFromOSRelease() *DistroInfo {
	if !FileExists("/etc/os-release") {
		return nil
	}

	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil
	}

	osRelease := string(content)
	return getDistroFromOSReleaseContent(osRelease)
}

type distroMapping struct {
	id        string
	getDistro func() *DistroInfo
}

// getDistroFromOSReleaseContent parses os-release content to identify distribution.
// osRelease: Content of /etc/os-release file. Returns matching DistroInfo or nil.
func getDistroFromOSReleaseContent(osRelease string) *DistroInfo {
	distros := []distroMapping{
		{"arch", getArchDistro},
		{"ubuntu", getUbuntuDistro},
		{"debian", getDebianDistro},
		{"fedora", getFedoraDistro},
		{"centos", getCentOSDistro},
		{"rhel", getRHELDistro},
		{"opensuse", getOpenSUSEDistro},
	}

	for _, distro := range distros {
		if containsDistroID(osRelease, distro.id) {
			return distro.getDistro()
		}
	}

	return nil
}

// containsDistroID checks if os-release content contains a specific distribution ID.
// content: os-release file content, id: Distribution identifier. Returns true if ID found.
func containsDistroID(content, id string) bool {
	return strings.Contains(content, "ID="+id) || strings.Contains(content, "ID=\""+id+"\"")
}

type packageManagerMapping struct {
	command   string
	getDistro func() *DistroInfo
}

// detectFromPackageManagers identifies distribution by checking available package managers.
// Returns DistroInfo based on found package manager or error if none recognized.
func detectFromPackageManagers() (*DistroInfo, error) {
	packageManagers := []packageManagerMapping{
		{"pacman", getArchBasedDistro},
		{aptGetCmd, getDebianBasedDistro},
		{"dnf", getFedoraBasedDistro},
		{"yum", getRHELBasedDistro},
		{"zypper", getSUSEBasedDistro},
	}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm.command); err == nil {
			return pm.getDistro(), nil
		}
	}

	return nil, fmt.Errorf("unsupported Linux distribution")
}

func createDistroInfo(name, packageManager string, installCmd, queryCmd []string) *DistroInfo {
	return &DistroInfo{
		Name:           name,
		PackageManager: packageManager,
		InstallCmd:     installCmd,
		QueryCmd:       queryCmd,
	}
}

func createArchDistro(name string) *DistroInfo {
	return createDistroInfo(name, "pacman",
		[]string{"pacman", "-S", "--needed", "--noconfirm"},
		[]string{"pacman", "-Q"})
}

func createDebianDistro(name string) *DistroInfo {
	return createDistroInfo(name, "apt",
		[]string{aptGetCmd, "install", "-y"},
		[]string{"dpkg", "-l"})
}

func createDNFDistro(name string) *DistroInfo {
	return createDistroInfo(name, "dnf",
		[]string{"dnf", "install", "-y"},
		[]string{"rpm", "-q"})
}

func createYumDistro(name string) *DistroInfo {
	return createDistroInfo(name, "yum",
		[]string{"yum", "install", "-y"},
		[]string{"rpm", "-q"})
}

func createZypperDistro(name string) *DistroInfo {
	return createDistroInfo(name, "zypper",
		[]string{"zypper", "install", "-y"},
		[]string{"rpm", "-q"})
}

func getArchDistro() *DistroInfo {
	return createArchDistro("Arch Linux")
}

func getUbuntuDistro() *DistroInfo {
	return createDebianDistro("Ubuntu")
}

func getDebianDistro() *DistroInfo {
	return createDebianDistro("Debian")
}

func getFedoraDistro() *DistroInfo {
	return createDNFDistro("Fedora")
}

func getCentOSDistro() *DistroInfo {
	return createYumDistro("CentOS")
}

func getRHELDistro() *DistroInfo {
	return createYumDistro("Red Hat Enterprise Linux")
}

func getOpenSUSEDistro() *DistroInfo {
	return createZypperDistro("openSUSE")
}

func getArchBasedDistro() *DistroInfo {
	return createArchDistro("Arch-based")
}

func getDebianBasedDistro() *DistroInfo {
	return createDebianDistro("Debian-based")
}

func getFedoraBasedDistro() *DistroInfo {
	return createDNFDistro("Fedora-based")
}

func getRHELBasedDistro() *DistroInfo {
	return createYumDistro("RHEL-based")
}

func getSUSEBasedDistro() *DistroInfo {
	return createZypperDistro("SUSE-based")
}

// GetBuildDependencies returns mapping of build tools to distribution-specific package names.
// Returns map where keys are generic tool names and values are PackageMapping structs.
func GetBuildDependencies() map[string]PackageMapping {
	return map[string]PackageMapping{
		"gcc": {
			Arch:   "gcc",
			Debian: "build-essential",
			RedHat: "gcc",
			SUSE:   "gcc",
		},
		"make": {
			Arch:   "make",
			Debian: "make",
			RedHat: "make",
			SUSE:   "make",
		},
		"autoconf": {
			Arch:   "autoconf",
			Debian: "autoconf",
			RedHat: "autoconf",
			SUSE:   "autoconf",
		},
		"pkgconf": {
			Arch:   "pkgconf",
			Debian: "pkg-config",
			RedHat: "pkgconfig",
			SUSE:   "pkg-config",
		},
		"bison": {
			Arch:   "bison",
			Debian: "bison",
			RedHat: "bison",
			SUSE:   "bison",
		},
		"re2c": {
			Arch:   "re2c",
			Debian: "re2c",
			RedHat: "re2c",
			SUSE:   "re2c",
		},
		"libxml2": {
			Arch:   "libxml2",
			Debian: "libxml2-dev",
			RedHat: "libxml2-devel",
			SUSE:   "libxml2-devel",
		},
		"curl": {
			Arch:   "curl",
			Debian: "libcurl4-openssl-dev",
			RedHat: "libcurl-devel",
			SUSE:   "libcurl-devel",
		},
		"openssl": {
			Arch:   "openssl",
			Debian: "libssl-dev",
			RedHat: "openssl-devel",
			SUSE:   "libopenssl-devel",
		},
		"zlib": {
			Arch:   "zlib",
			Debian: "zlib1g-dev",
			RedHat: "zlib-devel",
			SUSE:   "zlib-devel",
		},
		"libpng": {
			Arch:   "libpng",
			Debian: "libpng-dev",
			RedHat: "libpng-devel",
			SUSE:   "libpng16-devel",
		},
		"libjpeg-turbo": {
			Arch:   "libjpeg-turbo",
			Debian: "libjpeg-dev",
			RedHat: "libjpeg-turbo-devel",
			SUSE:   "libjpeg62-devel",
		},
		"freetype2": {
			Arch:   "freetype2",
			Debian: "libfreetype6-dev",
			RedHat: "freetype-devel",
			SUSE:   "freetype2-devel",
		},
	}
}

// GetPackageForDistro returns the correct package name for a tool on a specific distribution.
// packageKey: Generic package identifier, distro: Target distribution. Returns package name.
func GetPackageForDistro(packageKey string, distro *DistroInfo) string {
	packages := GetBuildDependencies()
	mapping, exists := packages[packageKey]
	if !exists {
		return packageKey
	}

	switch distro.PackageManager {
	case "pacman":
		return mapping.Arch
	case "apt":
		return mapping.Debian
	case "dnf", "yum":
		return mapping.RedHat
	case "zypper":
		return mapping.SUSE
	default:
		return packageKey
	}
}
