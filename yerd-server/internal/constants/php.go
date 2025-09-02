package constants

type Extension struct {
	Name         string
	ConfigFlag   string
	Dependencies []string
	PhpVersions  []string
	IsPECL       bool
	PECLName     string
}

var PhpVersions = []string{"8.1", "8.2", "8.3", "8.4"}
var DefaultExtensions = []string{
	"mbstring", "curl", "openssl", "fileinfo", "filter", "hash",
	"pcre", "session", "xml", "zip", "sqlite3", "sockets", "zlib",
}

var availableExtensions = map[string]Extension{
	"mbstring": {
		Name:        "mbstring",
		ConfigFlag:  "--enable-mbstring",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"bcmath": {
		Name:        "bcmath",
		ConfigFlag:  "--enable-bcmath",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"opcache": {
		Name:        "opcache",
		ConfigFlag:  "--enable-opcache",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"curl": {
		Name:         "curl",
		ConfigFlag:   "--with-curl",
		Dependencies: []string{"libcurl"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"openssl": {
		Name:         "openssl",
		ConfigFlag:   "--with-openssl",
		Dependencies: []string{"openssl"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"zip": {
		Name:         "zip",
		ConfigFlag:   "--with-zip",
		Dependencies: []string{"libzip"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"sockets": {
		Name:        "sockets",
		ConfigFlag:  "--enable-sockets",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"mysqli": {
		Name:         "mysqli",
		ConfigFlag:   "--with-mysqli",
		Dependencies: []string{"mysql"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"pdo-mysql": {
		Name:         "pdo-mysql",
		ConfigFlag:   "--with-pdo-mysql",
		Dependencies: []string{"mysql"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"gd": {
		Name:         "gd",
		ConfigFlag:   "--enable-gd",
		Dependencies: []string{"libgd"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"jpeg": {
		Name:         "jpeg",
		ConfigFlag:   "--with-jpeg",
		Dependencies: []string{"libjpeg"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"freetype": {
		Name:         "freetype",
		ConfigFlag:   "--with-freetype",
		Dependencies: []string{"freetype2"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"xml": {
		Name:        "xml",
		ConfigFlag:  "--enable-xml",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"json": {
		Name:        "json",
		ConfigFlag:  "--enable-json",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"session": {
		Name:        "session",
		ConfigFlag:  "--enable-session",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"hash": {
		Name:        "hash",
		ConfigFlag:  "--enable-hash",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"filter": {
		Name:        "filter",
		ConfigFlag:  "--enable-filter",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"pcre": {
		Name:         "pcre",
		ConfigFlag:   "--with-pcre-jit",
		Dependencies: []string{"pcre2"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"zlib": {
		Name:         "zlib",
		ConfigFlag:   "--with-zlib",
		Dependencies: []string{"zlib"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"bz2": {
		Name:         "bz2",
		ConfigFlag:   "--with-bz2",
		Dependencies: []string{"bzip2"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"iconv": {
		Name:        "iconv",
		ConfigFlag:  "--with-iconv",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"intl": {
		Name:         "intl",
		ConfigFlag:   "--enable-intl",
		Dependencies: []string{"icu"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"pgsql": {
		Name:         "pgsql",
		ConfigFlag:   "--with-pgsql",
		Dependencies: []string{"postgresql"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"pdo-pgsql": {
		Name:         "pdo-pgsql",
		ConfigFlag:   "--with-pdo-pgsql",
		Dependencies: []string{"postgresql"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"sqlite3": {
		Name:         "sqlite3",
		ConfigFlag:   "--with-sqlite3",
		Dependencies: []string{"sqlite"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"pdo-sqlite": {
		Name:         "pdo-sqlite",
		ConfigFlag:   "--with-pdo-sqlite",
		Dependencies: []string{"sqlite"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"fileinfo": {
		Name:        "fileinfo",
		ConfigFlag:  "--enable-fileinfo",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
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
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"gmp": {
		Name:         "gmp",
		ConfigFlag:   "--with-gmp",
		Dependencies: []string{"gmp"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"ldap": {
		Name:         "ldap",
		ConfigFlag:   "--with-ldap",
		Dependencies: []string{"ldap"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       false,
	},
	"soap": {
		Name:        "soap",
		ConfigFlag:  "--enable-soap",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"ftp": {
		Name:        "ftp",
		ConfigFlag:  "--enable-ftp",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"pcntl": {
		Name:        "pcntl",
		ConfigFlag:  "--enable-pcntl",
		PhpVersions: []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:      false,
	},
	"imagick": {
		Name:         "imagick",
		ConfigFlag:   "",
		Dependencies: []string{"imagick"},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       true,
		PECLName:     "imagick",
	},
	"imap": {
		Name:         "imap",
		ConfigFlag:   "",
		Dependencies: []string{"imap"},
		PhpVersions:  []string{"8.4"},
		IsPECL:       true,
		PECLName:     "imap",
	},
	"redis": {
		Name:         "redis",
		ConfigFlag:   "",
		Dependencies: []string{},
		PhpVersions:  []string{"8.1", "8.2", "8.3", "8.4"},
		IsPECL:       true,
		PECLName:     "redis",
	},
}

func GetExtensionConfigureFlags(extensions []string) []string {
	var flags []string

	for _, extName := range extensions {
		if ext, exists := availableExtensions[extName]; exists {
			flags = append(flags, ext.ConfigFlag)
		}
	}

	return flags
}

func GetExtension(name string) (Extension, bool) {
	ext, exists := availableExtensions[name]
	return ext, exists
}

func GetExtensionDependencies(extensions []string) []string {
	depMap := make(map[string]bool)

	for _, extName := range extensions {
		if ext, exists := availableExtensions[extName]; exists {
			for _, dep := range ext.Dependencies {
				depMap[dep] = true
			}
		}
	}

	dependencies := make([]string, 0, len(depMap))
	for dep := range depMap {
		dependencies = append(dependencies, dep)
	}

	return dependencies
}
