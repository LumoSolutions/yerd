package constants

import (
	"maps"
	"slices"
	"sort"
	"strings"
)

type Extension struct {
	Name         string
	ConfigFlag   string
	Dependencies []string
	IsPECL       bool
	PECLName     string
}

var availablePhpVersions = []string{"8.1", "8.2", "8.3", "8.4"}
var availableExtensions = map[string]Extension{
	"mbstring": {
		Name:       "mbstring",
		ConfigFlag: "--enable-mbstring",
		IsPECL:     false,
	},
	"bcmath": {
		Name:       "bcmath",
		ConfigFlag: "--enable-bcmath",
		IsPECL:     false,
	},
	"opcache": {
		Name:       "opcache",
		ConfigFlag: "--enable-opcache",
		IsPECL:     false,
	},
	"curl": {
		Name:         "curl",
		ConfigFlag:   "--with-curl",
		Dependencies: []string{"libcurl"},
		IsPECL:       false,
	},
	"openssl": {
		Name:         "openssl",
		ConfigFlag:   "--with-openssl",
		Dependencies: []string{"openssl"},
		IsPECL:       false,
	},
	"zip": {
		Name:         "zip",
		ConfigFlag:   "--with-zip",
		Dependencies: []string{"libzip"},
		IsPECL:       false,
	},
	"sockets": {
		Name:       "sockets",
		ConfigFlag: "--enable-sockets",
		IsPECL:     false,
	},
	"mysqli": {
		Name:         "mysqli",
		ConfigFlag:   "--with-mysqli",
		Dependencies: []string{"mysql"},
		IsPECL:       false,
	},
	"pdo-mysql": {
		Name:         "pdo-mysql",
		ConfigFlag:   "--with-pdo-mysql",
		Dependencies: []string{"mysql"},
		IsPECL:       false,
	},
	"gd": {
		Name:         "gd",
		ConfigFlag:   "--enable-gd",
		Dependencies: []string{"libgd"},
		IsPECL:       false,
	},
	"jpeg": {
		Name:         "jpeg",
		ConfigFlag:   "--with-jpeg",
		Dependencies: []string{"libjpeg"},
		IsPECL:       false,
	},
	"freetype": {
		Name:         "freetype",
		ConfigFlag:   "--with-freetype",
		Dependencies: []string{"freetype2"},
		IsPECL:       false,
	},
	"xml": {
		Name:       "xml",
		ConfigFlag: "--enable-xml",
		IsPECL:     false,
	},
	"json": {
		Name:       "json",
		ConfigFlag: "--enable-json",
		IsPECL:     false,
	},
	"session": {
		Name:       "session",
		ConfigFlag: "--enable-session",
		IsPECL:     false,
	},
	"hash": {
		Name:       "hash",
		ConfigFlag: "--enable-hash",
		IsPECL:     false,
	},
	"filter": {
		Name:       "filter",
		ConfigFlag: "--enable-filter",
		IsPECL:     false,
	},
	"pcre": {
		Name:         "pcre",
		ConfigFlag:   "--with-pcre-jit",
		Dependencies: []string{"pcre2"},
		IsPECL:       false,
	},
	"zlib": {
		Name:         "zlib",
		ConfigFlag:   "--with-zlib",
		Dependencies: []string{"zlib"},
		IsPECL:       false,
	},
	"bz2": {
		Name:         "bz2",
		ConfigFlag:   "--with-bz2",
		Dependencies: []string{"bzip2"},
		IsPECL:       false,
	},
	"iconv": {
		Name:       "iconv",
		ConfigFlag: "--with-iconv",
		IsPECL:     false,
	},
	"intl": {
		Name:         "intl",
		ConfigFlag:   "--enable-intl",
		Dependencies: []string{"icu"},
		IsPECL:       false,
	},
	"pgsql": {
		Name:         "pgsql",
		ConfigFlag:   "--with-pgsql",
		Dependencies: []string{"postgresql"},
		IsPECL:       false,
	},
	"pdo-pgsql": {
		Name:         "pdo-pgsql",
		ConfigFlag:   "--with-pdo-pgsql",
		Dependencies: []string{"postgresql"},
		IsPECL:       false,
	},
	"sqlite3": {
		Name:         "sqlite3",
		ConfigFlag:   "--with-sqlite3",
		Dependencies: []string{"sqlite"},
		IsPECL:       false,
	},
	"pdo-sqlite": {
		Name:         "pdo-sqlite",
		ConfigFlag:   "--with-pdo-sqlite",
		Dependencies: []string{"sqlite"},
		IsPECL:       false,
	},
	"fileinfo": {
		Name:       "fileinfo",
		ConfigFlag: "--enable-fileinfo",
		IsPECL:     false,
	},
	"exif": {
		Name:       "exif",
		ConfigFlag: "--enable-exif",
		IsPECL:     false,
	},
	"gettext": {
		Name:         "gettext",
		ConfigFlag:   "--with-gettext",
		Dependencies: []string{"gettext"},
		IsPECL:       false,
	},
	"gmp": {
		Name:         "gmp",
		ConfigFlag:   "--with-gmp",
		Dependencies: []string{"gmp"},
		IsPECL:       false,
	},
	"ldap": {
		Name:         "ldap",
		ConfigFlag:   "--with-ldap",
		Dependencies: []string{"ldap"},
		IsPECL:       false,
	},
	"soap": {
		Name:       "soap",
		ConfigFlag: "--enable-soap",
		IsPECL:     false,
	},
	"ftp": {
		Name:       "ftp",
		ConfigFlag: "--enable-ftp",
		IsPECL:     false,
	},
	"pcntl": {
		Name:       "pcntl",
		ConfigFlag: "--enable-pcntl",
		IsPECL:     false,
	},
	"imagick": {
		Name:         "imagick",
		ConfigFlag:   "",
		Dependencies: []string{"imagick"},
		IsPECL:       true,
		PECLName:     "imagick",
	},
	"imap": {
		Name:         "imap",
		ConfigFlag:   "--with-imap",
		Dependencies: []string{"imap"},
		IsPECL:       false,
	},
	"redis": {
		Name:         "redis",
		ConfigFlag:   "",
		Dependencies: []string{},
		IsPECL:       true,
		PECLName:     "redis",
	},
}
var defaultExtensions = []string{
	"mbstring", "curl", "openssl", "fileinfo", "filter", "hash",
	"pcre", "session", "xml", "zip", "mysqli", "sqlite3", "pdo-mysql",
	"sockets", "zlib",
}

// GetAvailableVersions returns the list of PHP versions supported by YERD.
func GetAvailablePhpVersions() []string {
	return availablePhpVersions
}

// IsValidVersion checks if the provided version string is supported by YERD.
// version: PHP version string to validate. Returns true if version is supported.
func IsValidPhpVersion(version string) bool {
	return slices.Contains(availablePhpVersions, version)
}

// GetDefaultExtensions returns the default extensions for all PHP installations
func GetDefaultExtensions() []string {
	return defaultExtensions
}

// GetExtension retrieves extension information by name.
// name: Extension name to lookup. Returns Extension struct and existence boolean.
func GetExtensions() []string {
	keys := make([]string, 0, len(availableExtensions))
	for key := range availableExtensions {
		keys = append(keys, key)
	}
	return keys
}

// GetExtension retrieves extension information by name.
// name: Extension name to lookup. Returns Extension struct and existence boolean.
func GetExtension(name string) (Extension, bool) {
	ext, exists := availableExtensions[name]
	return ext, exists
}

// ValidateExtensions separates provided extensions into valid and invalid lists.
// extensions: Extension names to validate. Returns valid extensions slice and invalid extensions slice.
func ValidateExtensions(extensions []string) ([]string, []string) {
	var valid []string
	var invalid []string

	for _, ext := range extensions {
		if _, exists := availableExtensions[ext]; exists {
			valid = append(valid, ext)
		} else {
			invalid = append(invalid, ext)
		}
	}

	return valid, invalid
}

// GetConfigureFlags returns PHP configure flags for the specified extensions.
// extensions: Extension names to get flags for. Returns slice of configure flag strings.
func GetExtensionConfigureFlags(extensions []string) []string {
	var flags []string

	for _, extName := range extensions {
		if ext, exists := availableExtensions[extName]; exists {
			flags = append(flags, ext.ConfigFlag)
		}
	}

	return flags
}

// GetDependencies returns system dependencies required for the specified extensions.
// extensions: Extension names to check. Returns sorted slice of unique system dependencies.
func GetExtensionDependencies(extensions []string) []string {
	depMap := make(map[string]bool)

	for _, extName := range extensions {
		if ext, exists := availableExtensions[extName]; exists {
			for _, dep := range ext.Dependencies {
				depMap[dep] = true
			}
		}
	}

	var deps []string
	for dep := range depMap {
		deps = append(deps, dep)
	}
	sort.Strings(deps)
	return deps
}

// SuggestSimilarExtensions finds extension names similar to an invalid extension name.
// invalid: Invalid extension name to find suggestions for. Returns sorted slice of similar extension names.
func SuggestSimilarExtensions(invalid string) []string {
	var suggestions []string
	invalid = strings.ToLower(invalid)

	for name := range availableExtensions {
		if strings.Contains(strings.ToLower(name), invalid) ||
			strings.Contains(invalid, strings.ToLower(name)) {
			suggestions = append(suggestions, name)
		}
	}

	sort.Strings(suggestions)
	return suggestions
}

// GetAvailableExtensions returns a list of available PHP extensions
func GetAvailableExtensions() []string {
	keys := maps.Keys(availableExtensions)
	return slices.Collect(keys)
}
