package config

type WebConfig struct {
	Installed bool
	Sites     map[string]SiteConfig
}

type SiteConfig struct {
	RootDirectory   string
	PublicDirectory string
	Domain          string
	PhpVersion      string
	NginxConfig     string
}
