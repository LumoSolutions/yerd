package config

type NginxConfig struct {
	Installed   bool
	Elevated    bool
	HttpPort    int
	HttpsPort   int
	ServiceName string
	Sites       map[string]SiteConfig
	CaCert      *CertInfo
}

type SiteConfig struct {
	RootDirectory   string
	PublicDirectory string
	Domain          string
	Aliases         []string
	Secured         bool
	PhpVersion      string
	Port            int
	SslPort         int
	SiteFile        string
	Certificate     *CertInfo
}
