package dependencies

// PackageManager represents different Linux package managers
type PackageManager string

const (
	APT    PackageManager = "apt"
	YUM    PackageManager = "yum"
	DNF    PackageManager = "dnf"
	PACMAN PackageManager = "pacman"
	ZYPPER PackageManager = "zypper"
	APKL   PackageManager = "apk"
)

// PackageManagerConfig defines command and arguments for each package manager
type PackageManagerConfig struct {
	Command     string
	CheckName   string
	InstallArgs []string
	QueryCmd    string
	QueryArgs   []string
}

// PackageManagerConfigs maps package managers to their configurations
var PackageManagerConfigs = map[PackageManager]PackageManagerConfig{
	APT: {
		Command:     "apt-get",
		CheckName:   "apt-get",
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    "dpkg",
		QueryArgs:   []string{"-l"},
	},
	YUM: {
		Command:     "yum",
		CheckName:   "yum",
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    "rpm",
		QueryArgs:   []string{"-q"},
	},
	DNF: {
		Command:     "dnf",
		CheckName:   "dnf",
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    "rpm",
		QueryArgs:   []string{"-q"},
	},
	PACMAN: {
		Command:     "pacman",
		CheckName:   "pacman",
		InstallArgs: []string{"-S", "--noconfirm"},
		QueryCmd:    "pacman",
		QueryArgs:   []string{"-Q"},
	},
	ZYPPER: {
		Command:     "zypper",
		CheckName:   "zypper",
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    "rpm",
		QueryArgs:   []string{"-q"},
	},
	APKL: {
		Command:     "apk",
		CheckName:   "apk",
		InstallArgs: []string{"add"},
		QueryCmd:    "apk",
		QueryArgs:   []string{"info"},
	},
}

// DependencyConfig represents all configuration for a single dependency
type DependencyConfig struct {
	Name            string
	SystemPackages  map[PackageManager][]string
	PkgConfigNames  map[PackageManager][]string
	CommonPkgConfig []string
	Commands        []string
	Libraries       []string
}

// DependencyRegistry is the single source of truth for all dependency configurations
var DependencyRegistry = map[string]*DependencyConfig{
	"curl": {
		Name: "curl",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libcurl4-openssl-dev"},
			YUM:    {"libcurl-devel"},
			DNF:    {"libcurl-devel"},
			PACMAN: {"curl"},
			ZYPPER: {"libcurl-devel"},
			APKL:   {"curl-dev"},
		},
		CommonPkgConfig: []string{"libcurl"},
	},
	"openssl": {
		Name: "openssl",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libssl-dev"},
			YUM:    {"openssl-devel"},
			DNF:    {"openssl-devel"},
			PACMAN: {"openssl"},
			ZYPPER: {"openssl-devel"},
			APKL:   {"openssl-dev"},
		},
		CommonPkgConfig: []string{"openssl"},
	},
	"zip": {
		Name: "zip",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libzip-dev"},
			YUM:    {"libzip-devel"},
			DNF:    {"libzip-devel"},
			PACMAN: {"libzip"},
			ZYPPER: {"libzip-devel"},
			APKL:   {"libzip-dev"},
		},
		CommonPkgConfig: []string{"libzip"},
	},
	"gd": {
		Name: "gd",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libgd-dev"},
			YUM:    {"gd-devel"},
			DNF:    {"gd-devel"},
			PACMAN: {"gd"},
			ZYPPER: {"gd-devel"},
			APKL:   {"gd-dev"},
		},
		CommonPkgConfig: []string{"gdlib"},
	},
	"mysqli": {
		Name: "mysqli",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libmysqlclient-dev"},
			YUM:    {"mysql-devel"},
			DNF:    {"mysql-devel"},
			PACMAN: {"mariadb-libs"},
			ZYPPER: {"libmysqlclient-devel"},
			APKL:   {"mysql-dev"},
		},
		Commands:  []string{"mysql_config"},
		Libraries: []string{"libmysqlclient", "libmariadb"},
	},
	"pdo-mysql": {
		Name: "pdo-mysql",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libmysqlclient-dev"},
			YUM:    {"mysql-devel"},
			DNF:    {"mysql-devel"},
			PACMAN: {"mariadb-libs"},
			ZYPPER: {"libmysqlclient-devel"},
			APKL:   {"mysql-dev"},
		},
		Commands:  []string{"mysql_config"},
		Libraries: []string{"libmysqlclient", "libmariadb"},
	},
	"jpeg": {
		Name: "jpeg",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libjpeg-dev"},
			YUM:    {"libjpeg-turbo-devel"},
			DNF:    {"libjpeg-turbo-devel"},
			PACMAN: {"libjpeg-turbo"},
			ZYPPER: {"libjpeg8-devel"},
			APKL:   {"libjpeg-turbo-dev"},
		},
		CommonPkgConfig: []string{"libjpeg"},
	},
	"freetype": {
		Name: "freetype",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libfreetype6-dev"},
			YUM:    {"freetype-devel"},
			DNF:    {"freetype-devel"},
			PACMAN: {"freetype2"},
			ZYPPER: {"freetype2-devel"},
			APKL:   {"freetype-dev"},
		},
		CommonPkgConfig: []string{"freetype2"},
	},
	"zlib": {
		Name: "zlib",
		SystemPackages: map[PackageManager][]string{
			APT:    {"zlib1g-dev"},
			YUM:    {"zlib-devel"},
			DNF:    {"zlib-devel"},
			PACMAN: {"zlib"},
			ZYPPER: {"zlib-devel"},
			APKL:   {"zlib-dev"},
		},
		CommonPkgConfig: []string{"zlib"},
	},
	"bz2": {
		Name: "bz2",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libbz2-dev"},
			YUM:    {"bzip2-devel"},
			DNF:    {"bzip2-devel"},
			PACMAN: {"bzip2"},
			ZYPPER: {"libbz2-devel"},
			APKL:   {"bzip2-dev"},
		},
		CommonPkgConfig: []string{"bzip2"},
	},
	"intl": {
		Name: "intl",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libicu-dev"},
			YUM:    {"libicu-devel"},
			DNF:    {"libicu-devel"},
			PACMAN: {"icu"},
			ZYPPER: {"libicu-devel"},
			APKL:   {"icu-dev"},
		},
		CommonPkgConfig: []string{"icu-uc", "icu-io"},
	},
	"gettext": {
		Name: "gettext",
		SystemPackages: map[PackageManager][]string{
			APT:    {"gettext"},
			YUM:    {"gettext-devel"},
			DNF:    {"gettext-devel"},
			PACMAN: {"gettext"},
			ZYPPER: {"gettext-tools"},
			APKL:   {"gettext-dev"},
		},
		Commands:  []string{"gettext"},
		Libraries: []string{"libintl"},
	},
	"gmp": {
		Name: "gmp",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libgmp-dev"},
			YUM:    {"gmp-devel"},
			DNF:    {"gmp-devel"},
			PACMAN: {"gmp"},
			ZYPPER: {"gmp-devel"},
			APKL:   {"gmp-dev"},
		},
		Libraries: []string{"libgmp"},
	},
	"mysql": {
		Name: "mysql",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libmysqlclient-dev"},
			YUM:    {"mysql-devel"},
			DNF:    {"mysql-devel"},
			PACMAN: {"mariadb-libs"},
			ZYPPER: {"libmysqlclient-devel"},
			APKL:   {"mysql-dev"},
		},
		Commands:  []string{"mysql_config"},
		Libraries: []string{"libmysqlclient"},
	},
	"pgsql": {
		Name: "pgsql",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libpq-dev"},
			YUM:    {"postgresql-devel"},
			DNF:    {"postgresql-devel"},
			PACMAN: {"postgresql-libs"},
			ZYPPER: {"postgresql-devel"},
			APKL:   {"postgresql-dev"},
		},
		Commands:  []string{"pg_config"},
		Libraries: []string{"libpq"},
	},
	"postgresql": {
		Name: "postgresql",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libpq-dev"},
			YUM:    {"postgresql-devel"},
			DNF:    {"postgresql-devel"},
			PACMAN: {"postgresql-libs"},
			ZYPPER: {"postgresql-devel"},
			APKL:   {"postgresql-dev"},
		},
		Commands:  []string{"pg_config"},
		Libraries: []string{"libpq"},
	},
	"pdo-pgsql": {
		Name: "pdo-pgsql",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libpq-dev"},
			YUM:    {"postgresql-devel"},
			DNF:    {"postgresql-devel"},
			PACMAN: {"postgresql-libs"},
			ZYPPER: {"postgresql-devel"},
			APKL:   {"postgresql-dev"},
		},
		Commands:  []string{"pg_config"},
		Libraries: []string{"libpq"},
	},
	"sqlite": {
		Name: "sqlite",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libsqlite3-dev"},
			YUM:    {"sqlite-devel"},
			DNF:    {"sqlite-devel"},
			PACMAN: {"sqlite"},
			ZYPPER: {"sqlite3-devel"},
			APKL:   {"sqlite-dev"},
		},
		Commands:  []string{"sqlite3"},
		Libraries: []string{"libsqlite3"},
	},
	"xml": {
		Name: "xml",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libxml2-dev"},
			YUM:    {"libxml2-devel"},
			DNF:    {"libxml2-devel"},
			PACMAN: {"libxml2"},
			ZYPPER: {"libxml2-devel"},
			APKL:   {"libxml2-dev"},
		},
		CommonPkgConfig: []string{"libxml-2.0"},
	},
	"pcre2": {
		Name: "pcre2",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libpcre2-dev"},
			YUM:    {"pcre2-devel"},
			DNF:    {"pcre2-devel"},
			PACMAN: {"pcre2"},
			ZYPPER: {"pcre2-devel"},
			APKL:   {"pcre2-dev"},
		},
		CommonPkgConfig: []string{"libpcre2-8"},
	},
	"ldap": {
		Name: "ldap",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libldap2-dev"},
			YUM:    {"openldap-devel"},
			DNF:    {"openldap-devel"},
			PACMAN: {"libldap"},
			ZYPPER: {"openldap2-devel"},
			APKL:   {"openldap-dev"},
		},
		Commands:  []string{"ldapsearch"},
		Libraries: []string{"libldap"},
	},
	"oniguruma": {
		Name: "oniguruma",
		SystemPackages: map[PackageManager][]string{
			APT:    {"libonig-dev"},
			YUM:    {"oniguruma-devel"},
			DNF:    {"oniguruma-devel"},
			PACMAN: {"oniguruma"},
			ZYPPER: {"libonig-devel"},
			APKL:   {"oniguruma-dev"},
		},
	},
	"re2c": {
		Name: "re2c",
		SystemPackages: map[PackageManager][]string{
			APT:    {"re2c"},
			YUM:    {"re2c"},
			DNF:    {"re2c"},
			PACMAN: {"re2c"},
			ZYPPER: {"re2c"},
			APKL:   {"re2c"},
		},
		Commands: []string{"re2c"},
	},
	"autoconf": {
		Name: "autoconf",
		SystemPackages: map[PackageManager][]string{
			APT:    {"autoconf"},
			YUM:    {"autoconf"},
			DNF:    {"autoconf"},
			PACMAN: {"autoconf"},
			ZYPPER: {"autoconf"},
			APKL:   {"autoconf"},
		},
		Commands: []string{"autoconf"},
	},
	"pkgconfig": {
		Name: "pkgconfig",
		SystemPackages: map[PackageManager][]string{
			APT:    {"pkg-config"},
			YUM:    {"pkgconfig"},
			DNF:    {"pkgconf"},
			PACMAN: {"pkgconf"},
			ZYPPER: {"pkg-config"},
			APKL:   {"pkgconf"},
		},
		Commands: []string{"pkg-config"},
	},
	"buildtools": {
		Name: "buildtools",
		SystemPackages: map[PackageManager][]string{
			APT:    {"build-essential"},
			YUM:    {"gcc", "gcc-c++", "make"},
			DNF:    {"gcc", "gcc-c++", "make"},
			PACMAN: {"base-devel"},
			ZYPPER: {"gcc", "gcc-c++", "make"},
			APKL:   {"build-base"},
		},
	},
	"webtools": {
		Name: "webtools",
		SystemPackages: map[PackageManager][]string{
			APT:    {"wget", "tar"},
			YUM:    {"wget", "tar"},
			DNF:    {"wget", "tar"},
			PACMAN: {"wget", "tar"},
			ZYPPER: {"wget", "tar"},
			APKL:   {"wget", "tar"},
		},
		Commands: []string{"wget", "tar"},
	},
}

// GetDependencyConfig returns the dependency configuration for a given name
func GetDependencyConfig(name string) (*DependencyConfig, bool) {
	config, exists := DependencyRegistry[name]
	return config, exists
}

// GetSystemPackages returns the system packages for a dependency on a specific distro
func GetSystemPackages(depName string, pm PackageManager) ([]string, bool) {
	if config, exists := DependencyRegistry[depName]; exists {
		if packages, hasDistro := config.SystemPackages[pm]; hasDistro {
			return packages, true
		}
	}
	return nil, false
}

// GetAllDependencyNames returns all available dependency names
func GetAllDependencyNames() []string {
	names := make([]string, 0, len(DependencyRegistry))
	for name := range DependencyRegistry {
		names = append(names, name)
	}
	return names
}

// GetBuildDependencies returns build dependencies for PHP compilation
func GetBuildDependencies(pm PackageManager) []string {
	var packages []string

	// Core build tools
	buildDeps := []string{"buildtools", "autoconf", "pkgconfig", "re2c", "oniguruma", "xml", "sqlite"}

	for _, dep := range buildDeps {
		if pkgs, exists := GetSystemPackages(dep, pm); exists {
			packages = append(packages, pkgs...)
		}
	}

	return packages
}

// GetWebBuildDependencies returns build dependencies for web services compilation
func GetWebBuildDependencies(pm PackageManager) []string {
	var packages []string

	// Web build tools
	webDeps := []string{"buildtools", "webtools"}

	for _, dep := range webDeps {
		if pkgs, exists := GetSystemPackages(dep, pm); exists {
			packages = append(packages, pkgs...)
		}
	}

	return packages
}

// GetPackageManagerConfig returns the configuration for a package manager
func GetPackageManagerConfig(pm PackageManager) (PackageManagerConfig, bool) {
	config, exists := PackageManagerConfigs[pm]
	return config, exists
}

// GetPackageManagerCommand returns the command string for a package manager
func GetPackageManagerCommand(pm PackageManager) string {
	if config, exists := PackageManagerConfigs[pm]; exists {
		return config.Command
	}
	return ""
}

// GetPackageManagerCheckName returns the check name for package manager detection
func GetPackageManagerCheckName(pm PackageManager) string {
	if config, exists := PackageManagerConfigs[pm]; exists {
		return config.CheckName
	}
	return ""
}
