package config

type NginxConfig struct {
	Installed     bool                  `json:"is_installed"`
	Elevated      bool                  `json:"is_elevated"`
	ServiceName   string                `json:"service_name"`
	IsUserService bool                  `json:"is_user_service"`
	Sites         map[string]SiteConfig `json:"sites"`
	CaCert        CertInfo              `json:"ca_cert"`
}

type SiteConfig struct {
	RootDirectory   string   `json:"root_dir"`
	PublicDirectory string   `json:"public_dir"`
	Domain          string   `json:"domain"`
	Aliases         []string `json:"aliases"`
	Secured         bool     `json:"secured"`
	PhpVersion      string   `json:"php_version"`
	Port            int      `json:"port"`
	SslPort         int      `json:"ssl_port"`
	SiteFile        string   `json:"siteFile"`
	Certificate     CertInfo `json:"cert_info"`
}

// GetNginxConfig returns a pointer to the core nginx configuration,
// or a new instances where one does not exist
func GetNginxConfig() *NginxConfig {
	return &NginxConfig{
		Installed:     false,
		Elevated:      false,
		ServiceName:   "yerd-nginx",
		IsUserService: true,
		Sites:         make(map[string]SiteConfig),
	}
}
