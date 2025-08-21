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
}

var availablePhpVersions = []string{"8.1", "8.2", "8.3", "8.4"}
var availableExtensions = map[string]Extension{
	"mbstring": {
		Name:       "mbstring",
		ConfigFlag: "--enable-mbstring",
	},
	"bcmath": {
		Name:       "bcmath",
		ConfigFlag: "--enable-bcmath",
	},
	"opcache": {
		Name:       "opcache",
		ConfigFlag: "--enable-opcache",
	},
	"curl": {
		Name:         "curl",
		ConfigFlag:   "--with-curl",
		Dependencies: []string{"libcurl"},
	},
	"openssl": {
		Name:         "openssl",
		ConfigFlag:   "--with-openssl",
		Dependencies: []string{"openssl"},
	},
	"zip": {
		Name:         "zip",
		ConfigFlag:   "--with-zip",
		Dependencies: []string{"libzip"},
	},
	"sockets": {
		Name:       "sockets",
		ConfigFlag: "--enable-sockets",
	},
	"mysqli": {
		Name:         "mysqli",
		ConfigFlag:   "--with-mysqli",
		Dependencies: []string{"mysql"},
	},
	"pdo-mysql": {
		Name:         "pdo-mysql",
		ConfigFlag:   "--with-pdo-mysql",
		Dependencies: []string{"mysql"},
	},
	"gd": {
		Name:         "gd",
		ConfigFlag:   "--enable-gd",
		Dependencies: []string{"libgd"},
	},
	"jpeg": {
		Name:         "jpeg",
		ConfigFlag:   "--with-jpeg",
		Dependencies: []string{"libjpeg"},
	},
	"freetype": {
		Name:         "freetype",
		ConfigFlag:   "--with-freetype",
		Dependencies: []string{"freetype2"},
	},
	"xml": {
		Name:       "xml",
		ConfigFlag: "--enable-xml",
	},
	"json": {
		Name:       "json",
		ConfigFlag: "--enable-json",
	},
	"session": {
		Name:       "session",
		ConfigFlag: "--enable-session",
	},
	"hash": {
		Name:       "hash",
		ConfigFlag: "--enable-hash",
	},
	"filter": {
		Name:       "filter",
		ConfigFlag: "--enable-filter",
	},
	"pcre": {
		Name:         "pcre",
		ConfigFlag:   "--with-pcre-jit",
		Dependencies: []string{"pcre2"},
	},
	"zlib": {
		Name:         "zlib",
		ConfigFlag:   "--with-zlib",
		Dependencies: []string{"zlib"},
	},
	"bz2": {
		Name:         "bz2",
		ConfigFlag:   "--with-bz2",
		Dependencies: []string{"bzip2"},
	},
	"iconv": {
		Name:       "iconv",
		ConfigFlag: "--with-iconv",
	},
	"intl": {
		Name:         "intl",
		ConfigFlag:   "--enable-intl",
		Dependencies: []string{"icu"},
	},
	"pgsql": {
		Name:         "pgsql",
		ConfigFlag:   "--with-pgsql",
		Dependencies: []string{"postgresql"},
	},
	"pdo-pgsql": {
		Name:         "pdo-pgsql",
		ConfigFlag:   "--with-pdo-pgsql",
		Dependencies: []string{"postgresql"},
	},
	"sqlite3": {
		Name:         "sqlite3",
		ConfigFlag:   "--with-sqlite3",
		Dependencies: []string{"sqlite"},
	},
	"pdo-sqlite": {
		Name:         "pdo-sqlite",
		ConfigFlag:   "--with-pdo-sqlite",
		Dependencies: []string{"sqlite"},
	},
	"fileinfo": {
		Name:       "fileinfo",
		ConfigFlag: "--enable-fileinfo",
	},
	"exif": {
		Name:       "exif",
		ConfigFlag: "--enable-exif",
	},
	"gettext": {
		Name:         "gettext",
		ConfigFlag:   "--with-gettext",
		Dependencies: []string{"gettext"},
	},
	"gmp": {
		Name:         "gmp",
		ConfigFlag:   "--with-gmp",
		Dependencies: []string{"gmp"},
	},
	"ldap": {
		Name:         "ldap",
		ConfigFlag:   "--with-ldap",
		Dependencies: []string{"ldap"},
	},
	"soap": {
		Name:       "soap",
		ConfigFlag: "--enable-soap",
	},
	"ftp": {
		Name:       "ftp",
		ConfigFlag: "--enable-ftp",
	},
	"pcntl": {
		Name:       "pcntl",
		ConfigFlag: "--enable-pcntl",
	},
}
var defaultExtensions = []string{
	"mbstring", "curl", "openssl", "fileinfo", "filter", "hash",
	"pcre", "session", "xml", "zip", "mysqli", "sqlite3", "pdo-mysql",
	"sockets",
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
