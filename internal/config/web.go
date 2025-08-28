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

func GetWebConfig() *WebConfig {
	var webConfig *WebConfig
	err := GetStruct("web", &webConfig)
	if err != nil || webConfig == nil {
		webConfig = &WebConfig{
			Installed: false,
			Sites:     make(map[string]SiteConfig),
		}
	}

	return webConfig
}
