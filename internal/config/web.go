package config

type WebConfig struct {
	Installed bool                  `json:"is_installed"`
	Sites     map[string]SiteConfig `json:"sites"`
}

type SiteConfig struct {
	RootDirectory   string `json:"rootDir"`
	PublicDirectory string `json:"publicDir"`
	Domain          string `json:"domain"`
	PhpVersion      string `json:"php_version"`
}
