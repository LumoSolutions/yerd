package extensions

import (
	"fmt"
	"sort"
	"strings"
)

type Extension struct {
	Name         string
	ConfigFlag   string
	Description  string
	Dependencies []string
	Category     string
}

var AvailableExtensions = map[string]Extension{
	"mbstring": {
		Name:        "mbstring",
		ConfigFlag:  "--enable-mbstring",
		Description: "Multibyte String Functions",
		Category:    "core",
	},
	"bcmath": {
		Name:        "bcmath",
		ConfigFlag:  "--enable-bcmath",
		Description: "BC Math Functions",
		Category:    "math",
	},
	"opcache": {
		Name:        "opcache",
		ConfigFlag:  "--enable-opcache",
		Description: "OPCache for improved performance",
		Category:    "performance",
	},
	"curl": {
		Name:        "curl",
		ConfigFlag:  "--with-curl",
		Description: "Client URL Library",
		Dependencies: []string{"libcurl"},
		Category:    "network",
	},
	"openssl": {
		Name:        "openssl",
		ConfigFlag:  "--with-openssl",
		Description: "OpenSSL support",
		Dependencies: []string{"openssl"},
		Category:    "security",
	},
	"zip": {
		Name:        "zip",
		ConfigFlag:  "--with-zip",
		Description: "Zip Archive support",
		Dependencies: []string{"libzip"},
		Category:    "archive",
	},
	"sockets": {
		Name:        "sockets",
		ConfigFlag:  "--enable-sockets",
		Description: "Socket Functions",
		Category:    "network",
	},
	"mysqli": {
		Name:        "mysqli",
		ConfigFlag:  "--with-mysqli",
		Description: "MySQL Improved Extension",
		Dependencies: []string{"mysql"},
		Category:    "database",
	},
	"pdo-mysql": {
		Name:        "pdo-mysql",
		ConfigFlag:  "--with-pdo-mysql",
		Description: "PDO MySQL Driver",
		Dependencies: []string{"mysql"},
		Category:    "database",
	},
	"gd": {
		Name:        "gd",
		ConfigFlag:  "--enable-gd",
		Description: "GD Graphics Library",
		Dependencies: []string{"libgd"},
		Category:    "graphics",
	},
	"jpeg": {
		Name:        "jpeg",
		ConfigFlag:  "--with-jpeg",
		Description: "JPEG support for GD",
		Dependencies: []string{"libjpeg"},
		Category:    "graphics",
	},
	"freetype": {
		Name:        "freetype",
		ConfigFlag:  "--with-freetype",
		Description: "FreeType 2 support",
		Dependencies: []string{"freetype2"},
		Category:    "graphics",
	},
	"xml": {
		Name:        "xml",
		ConfigFlag:  "--enable-xml",
		Description: "XML Parser support",
		Category:    "data",
	},
	"json": {
		Name:        "json",
		ConfigFlag:  "--enable-json",
		Description: "JSON support",
		Category:    "data",
	},
	"session": {
		Name:        "session",
		ConfigFlag:  "--enable-session",
		Description: "Session support",
		Category:    "core",
	},
	"hash": {
		Name:        "hash",
		ConfigFlag:  "--enable-hash",
		Description: "HASH Message Digest Framework",
		Category:    "security",
	},
	"filter": {
		Name:        "filter",
		ConfigFlag:  "--enable-filter",
		Description: "Input Filter support",
		Category:    "security",
	},
	"pcre": {
		Name:        "pcre",
		ConfigFlag:  "--with-pcre-jit",
		Description: "Perl Compatible Regular Expressions with JIT",
		Dependencies: []string{"pcre2"},
		Category:    "core",
	},
	"zlib": {
		Name:        "zlib",
		ConfigFlag:  "--with-zlib",
		Description: "Zlib compression support",
		Dependencies: []string{"zlib"},
		Category:    "compression",
	},
	"bz2": {
		Name:        "bz2",
		ConfigFlag:  "--with-bz2",
		Description: "Bzip2 compression support",
		Dependencies: []string{"bzip2"},
		Category:    "compression",
	},
	"iconv": {
		Name:        "iconv",
		ConfigFlag:  "--with-iconv",
		Description: "Character set conversion support",
		Category:    "core",
	},
	"intl": {
		Name:        "intl",
		ConfigFlag:  "--enable-intl",
		Description: "Internationalization extension",
		Dependencies: []string{"icu"},
		Category:    "i18n",
	},
	"pgsql": {
		Name:        "pgsql",
		ConfigFlag:  "--with-pgsql",
		Description: "PostgreSQL support",
		Dependencies: []string{"postgresql"},
		Category:    "database",
	},
	"pdo-pgsql": {
		Name:        "pdo-pgsql",
		ConfigFlag:  "--with-pdo-pgsql",
		Description: "PDO PostgreSQL Driver",
		Dependencies: []string{"postgresql"},
		Category:    "database",
	},
	"sqlite3": {
		Name:        "sqlite3",
		ConfigFlag:  "--with-sqlite3",
		Description: "SQLite 3 support",
		Dependencies: []string{"sqlite"},
		Category:    "database",
	},
	"pdo-sqlite": {
		Name:        "pdo-sqlite",
		ConfigFlag:  "--with-pdo-sqlite",
		Description: "PDO SQLite Driver",
		Dependencies: []string{"sqlite"},
		Category:    "database",
	},
	"fileinfo": {
		Name:        "fileinfo",
		ConfigFlag:  "--enable-fileinfo",
		Description: "File Information support",
		Category:    "filesystem",
	},
	"exif": {
		Name:        "exif",
		ConfigFlag:  "--enable-exif",
		Description: "EXIF image information support",
		Category:    "graphics",
	},
	"gettext": {
		Name:        "gettext",
		ConfigFlag:  "--with-gettext",
		Description: "GNU gettext support",
		Dependencies: []string{"gettext"},
		Category:    "i18n",
	},
	"gmp": {
		Name:        "gmp",
		ConfigFlag:  "--with-gmp",
		Description: "GNU MP support",
		Dependencies: []string{"gmp"},
		Category:    "math",
	},
	"ldap": {
		Name:        "ldap",
		ConfigFlag:  "--with-ldap",
		Description: "LDAP support",
		Dependencies: []string{"openldap"},
		Category:    "directory",
	},
	"soap": {
		Name:        "soap",
		ConfigFlag:  "--enable-soap",
		Description: "SOAP support",
		Category:    "webservices",
	},
	"ftp": {
		Name:        "ftp",
		ConfigFlag:  "--enable-ftp",
		Description: "FTP support",
		Category:    "network",
	},
}

func GetExtension(name string) (Extension, bool) {
	ext, exists := AvailableExtensions[name]
	return ext, exists
}

func GetExtensionsByCategory(category string) []Extension {
	var extensions []Extension
	for _, ext := range AvailableExtensions {
		if ext.Category == category {
			extensions = append(extensions, ext)
		}
	}
	return extensions
}

func GetAllCategories() []string {
	categories := make(map[string]bool)
	for _, ext := range AvailableExtensions {
		categories[ext.Category] = true
	}
	
	var result []string
	for category := range categories {
		result = append(result, category)
	}
	sort.Strings(result)
	return result
}

func ValidateExtensions(extensions []string) ([]string, []string) {
	var valid []string
	var invalid []string
	
	for _, ext := range extensions {
		if _, exists := AvailableExtensions[ext]; exists {
			valid = append(valid, ext)
		} else {
			invalid = append(invalid, ext)
		}
	}
	
	return valid, invalid
}

func GetConfigureFlags(extensions []string) []string {
	var flags []string
	
	for _, extName := range extensions {
		if ext, exists := AvailableExtensions[extName]; exists {
			flags = append(flags, ext.ConfigFlag)
		}
	}
	
	return flags
}

func GetDependencies(extensions []string) []string {
	depMap := make(map[string]bool)
	
	for _, extName := range extensions {
		if ext, exists := AvailableExtensions[extName]; exists {
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

func FormatExtensionList(extensions []string, withDescriptions bool) string {
	if len(extensions) == 0 {
		return "None"
	}
	
	if !withDescriptions {
		return strings.Join(extensions, ", ")
	}
	
	var formatted []string
	for _, extName := range extensions {
		if ext, exists := AvailableExtensions[extName]; exists {
			formatted = append(formatted, fmt.Sprintf("%s (%s)", extName, ext.Description))
		} else {
			formatted = append(formatted, extName)
		}
	}
	
	return strings.Join(formatted, "\n")
}

func SuggestSimilarExtensions(invalid string) []string {
	var suggestions []string
	invalid = strings.ToLower(invalid)
	
	for name := range AvailableExtensions {
		if strings.Contains(strings.ToLower(name), invalid) || 
		   strings.Contains(invalid, strings.ToLower(name)) {
			suggestions = append(suggestions, name)
		}
	}
	
	sort.Strings(suggestions)
	return suggestions
}