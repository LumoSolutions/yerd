package constants

const (
	APT            string = "apt"
	APTGET         string = "apt-get"
	YUM            string = "yum"
	DNF            string = "dnf"
	PACMAN         string = "pacman"
	ZYPPER         string = "zypper"
	APKL           string = "apk"
	RPM            string = "rpm"
	LIBCURLDEVEL   string = "libcurl-devel"
	OPENSSLDEVEL   string = "openssl-devel"
	LIBZIPDEVEL    string = "libzip-devel"
	GDDEVEL        string = "gd-devel"
	LIBMYSQLCLIENT string = "libmysqlclient-dev"
	MYSQLDEVEL     string = "mysql-devel"
	MARIADBLIBS    string = "mariadb-libs"
	LIBSMYSQLDEVEL string = "libmysqlclient-devel"
	MYSQLDEV       string = "mysql-dev"
	ZLIBDEVEL      string = "zlib-devel"
	LIBICUDEVEL    string = "libicu-devel"
	GMPDEVEL       string = "gmp-devel"
	LIBPQDEV       string = "libpq-dev"
	POSTGRESDEVEL  string = "postgresql-devel"
	POSTGRESLIBS   string = "postgresql-libs"
	POSTGRESDEV    string = "postgresql-dev"
	LIBXML2DEVEL   string = "libxml2-devel"
	PCRE2DEVEL     string = "pcre2-devel"
	PKGCONFIG      string = "pkg-config"
	GCCC           string = "gcc-c++"
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
var PackageManagerConfigs = map[string]PackageManagerConfig{
	APT: {
		Command:     APTGET,
		CheckName:   APTGET,
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    "dpkg",
		QueryArgs:   []string{"-l"},
	},
	YUM: {
		Command:     YUM,
		CheckName:   YUM,
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    RPM,
		QueryArgs:   []string{"-q"},
	},
	DNF: {
		Command:     DNF,
		CheckName:   DNF,
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    RPM,
		QueryArgs:   []string{"-q"},
	},
	PACMAN: {
		Command:     PACMAN,
		CheckName:   PACMAN,
		InstallArgs: []string{"-S", "--noconfirm"},
		QueryCmd:    PACMAN,
		QueryArgs:   []string{"-Q"},
	},
	ZYPPER: {
		Command:     ZYPPER,
		CheckName:   ZYPPER,
		InstallArgs: []string{"install", "-y"},
		QueryCmd:    RPM,
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
	SystemPackages  map[string][]string
	PkgConfigNames  map[string][]string
	CommonPkgConfig []string
	Commands        []string
	Libraries       []string
}

// DependencyRegistry is the single source of truth for all dependency configurations
var DependencyRegistry = map[string]*DependencyConfig{
	"curl": {
		Name: "curl",
		SystemPackages: map[string][]string{
			APT:    {"libcurl4-openssl-dev"},
			YUM:    {LIBCURLDEVEL},
			DNF:    {LIBCURLDEVEL},
			PACMAN: {"curl"},
			ZYPPER: {LIBCURLDEVEL},
			APKL:   {"curl-dev"},
		},
		CommonPkgConfig: []string{"libcurl"},
	},
	"openssl": {
		Name: "openssl",
		SystemPackages: map[string][]string{
			APT:    {"libssl-dev"},
			YUM:    {OPENSSLDEVEL},
			DNF:    {OPENSSLDEVEL},
			PACMAN: {"openssl"},
			ZYPPER: {OPENSSLDEVEL},
			APKL:   {"openssl-dev"},
		},
		CommonPkgConfig: []string{"openssl"},
	},
	"zip": {
		Name: "zip",
		SystemPackages: map[string][]string{
			APT:    {"libzip-dev"},
			YUM:    {LIBZIPDEVEL},
			DNF:    {LIBZIPDEVEL},
			PACMAN: {"libzip"},
			ZYPPER: {LIBZIPDEVEL},
			APKL:   {"libzip-dev"},
		},
		CommonPkgConfig: []string{"libzip"},
	},
	"gd": {
		Name: "gd",
		SystemPackages: map[string][]string{
			APT:    {"libgd-dev"},
			YUM:    {GDDEVEL},
			DNF:    {GDDEVEL},
			PACMAN: {"gd"},
			ZYPPER: {GDDEVEL},
			APKL:   {"gd-dev"},
		},
		CommonPkgConfig: []string{"gdlib"},
	},
	"mysqli": {
		Name: "mysqli",
		SystemPackages: map[string][]string{
			APT:    {LIBMYSQLCLIENT},
			YUM:    {MYSQLDEVEL},
			DNF:    {MYSQLDEVEL},
			PACMAN: {MARIADBLIBS},
			ZYPPER: {LIBSMYSQLDEVEL},
			APKL:   {MYSQLDEV},
		},
		Commands:  []string{"mysql_config"},
		Libraries: []string{"libmysqlclient", "libmariadb"},
	},
	"pdo-mysql": {
		Name: "pdo-mysql",
		SystemPackages: map[string][]string{
			APT:    {LIBMYSQLCLIENT},
			YUM:    {MYSQLDEVEL},
			DNF:    {MYSQLDEVEL},
			PACMAN: {MARIADBLIBS},
			ZYPPER: {LIBSMYSQLDEVEL},
			APKL:   {MYSQLDEV},
		},
		Commands:  []string{"mysql_config"},
		Libraries: []string{"libmysqlclient", "libmariadb"},
	},
	"jpeg": {
		Name: "jpeg",
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
			APT:    {"zlib1g-dev"},
			YUM:    {ZLIBDEVEL},
			DNF:    {ZLIBDEVEL},
			PACMAN: {"zlib"},
			ZYPPER: {ZLIBDEVEL},
			APKL:   {"zlib-dev"},
		},
		CommonPkgConfig: []string{"zlib"},
	},
	"bz2": {
		Name: "bz2",
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
			APT:    {"libicu-dev"},
			YUM:    {LIBICUDEVEL},
			DNF:    {LIBICUDEVEL},
			PACMAN: {"icu"},
			ZYPPER: {LIBICUDEVEL},
			APKL:   {"icu-dev"},
		},
		CommonPkgConfig: []string{"icu-uc", "icu-io"},
	},
	"gettext": {
		Name: "gettext",
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
			APT:    {"libgmp-dev"},
			YUM:    {GMPDEVEL},
			DNF:    {GMPDEVEL},
			PACMAN: {"gmp"},
			ZYPPER: {GMPDEVEL},
			APKL:   {"gmp-dev"},
		},
		Libraries: []string{"libgmp"},
	},
	"mysql": {
		Name: "mysql",
		SystemPackages: map[string][]string{
			APT:    {LIBMYSQLCLIENT},
			YUM:    {MYSQLDEVEL},
			DNF:    {MYSQLDEVEL},
			PACMAN: {MARIADBLIBS},
			ZYPPER: {LIBSMYSQLDEVEL},
			APKL:   {MYSQLDEV},
		},
		Commands:  []string{"mysql_config"},
		Libraries: []string{"libmysqlclient"},
	},
	"pgsql": {
		Name: "pgsql",
		SystemPackages: map[string][]string{
			APT:    {LIBPQDEV},
			YUM:    {POSTGRESDEVEL},
			DNF:    {POSTGRESDEVEL},
			PACMAN: {POSTGRESLIBS},
			ZYPPER: {POSTGRESDEVEL},
			APKL:   {POSTGRESDEV},
		},
		Commands:  []string{"pg_config"},
		Libraries: []string{"libpq"},
	},
	"postgresql": {
		Name: "postgresql",
		SystemPackages: map[string][]string{
			APT:    {LIBPQDEV},
			YUM:    {POSTGRESDEVEL},
			DNF:    {POSTGRESDEVEL},
			PACMAN: {POSTGRESLIBS},
			ZYPPER: {POSTGRESDEVEL},
			APKL:   {POSTGRESDEV},
		},
		Commands:  []string{"pg_config"},
		Libraries: []string{"libpq"},
	},
	"pdo-pgsql": {
		Name: "pdo-pgsql",
		SystemPackages: map[string][]string{
			APT:    {LIBPQDEV},
			YUM:    {POSTGRESDEVEL},
			DNF:    {POSTGRESDEVEL},
			PACMAN: {POSTGRESLIBS},
			ZYPPER: {POSTGRESDEVEL},
			APKL:   {POSTGRESDEV},
		},
		Commands:  []string{"pg_config"},
		Libraries: []string{"libpq"},
	},
	"sqlite": {
		Name: "sqlite",
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
			APT:    {"libxml2-dev"},
			YUM:    {LIBXML2DEVEL},
			DNF:    {LIBXML2DEVEL},
			PACMAN: {"libxml2"},
			ZYPPER: {LIBXML2DEVEL},
			APKL:   {"libxml2-dev"},
		},
		CommonPkgConfig: []string{"libxml-2.0"},
	},
	"pcre2": {
		Name: "pcre2",
		SystemPackages: map[string][]string{
			APT:    {"libpcre2-dev"},
			YUM:    {PCRE2DEVEL},
			DNF:    {PCRE2DEVEL},
			PACMAN: {"pcre2"},
			ZYPPER: {PCRE2DEVEL},
			APKL:   {"pcre2-dev"},
		},
		CommonPkgConfig: []string{"libpcre2-8"},
	},
	"ldap": {
		Name: "ldap",
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
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
		SystemPackages: map[string][]string{
			APT:    {PKGCONFIG},
			YUM:    {"pkgconfig"},
			DNF:    {"pkgconf"},
			PACMAN: {"pkgconf"},
			ZYPPER: {PKGCONFIG},
			APKL:   {"pkgconf"},
		},
		Commands: []string{PKGCONFIG},
	},
	"buildtools": {
		Name: "buildtools",
		SystemPackages: map[string][]string{
			APT:    {"build-essential"},
			YUM:    {"gcc", GCCC, "make"},
			DNF:    {"gcc", GCCC, "make"},
			PACMAN: {"base-devel"},
			ZYPPER: {"gcc", GCCC, "make"},
			APKL:   {"build-base"},
		},
	},
	"webtools": {
		Name: "webtools",
		SystemPackages: map[string][]string{
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
func GetSystemPackages(depName string, pm string) ([]string, bool) {
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
func GetBuildDependencies(pm string) []string {
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
func GetWebBuildDependencies(pm string) []string {
	var packages []string

	// Web build tools
	nginxConfig := GetNginxConfig()
	webDeps := []string{"buildtools", "webtools"}
	webDeps = append(webDeps, nginxConfig.Dependencies...)

	for _, dep := range webDeps {
		if pkgs, exists := GetSystemPackages(dep, pm); exists {
			packages = append(packages, pkgs...)
		}
	}

	return packages
}

// GetPackageManagerConfig returns the configuration for a package manager
func GetPackageManagerConfig(pm string) (PackageManagerConfig, bool) {
	config, exists := PackageManagerConfigs[pm]
	return config, exists
}

// GetPackageManagerCommand returns the command string for a package manager
func GetPackageManagerCommand(pm string) string {
	if config, exists := PackageManagerConfigs[pm]; exists {
		return config.Command
	}
	return ""
}

// GetPackageManagerCheckName returns the check name for package manager detection
func GetPackageManagerCheckName(pm string) string {
	if config, exists := PackageManagerConfigs[pm]; exists {
		return config.CheckName
	}
	return ""
}
