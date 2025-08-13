package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	YerdHome       string              `json:"yerd_home"`
	InstalledPHP   map[string]PHPInfo  `json:"installed_php"`
	CurrentCLI     string              `json:"current_cli"`
	PathConfigured bool                `json:"path_configured"`
}

type PHPInfo struct {
	Version     string    `json:"version"`
	InstallPath string    `json:"install_path"`
	InstallDate time.Time `json:"install_date"`
	IsCLI       bool      `json:"is_cli"`
	Extensions  []string  `json:"extensions"`
}

func GetConfigPath() (string, error) {
	var homeDir string
	

	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {

		realUser, err := user.Lookup(sudoUser)
		if err != nil {
			return "", err
		}
		homeDir = realUser.HomeDir
	} else {

		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	
	configDir := filepath.Join(homeDir, ".config", "yerd")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return getDefaultConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	
	if config.InstalledPHP == nil {
		config.InstalledPHP = make(map[string]PHPInfo)
	}
	
	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	
	if err := c.writeConfigFile(configPath); err != nil {
		return err
	}
	
	return fixFileOwnership(configPath)
}

func (c *Config) writeConfigFile(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(configPath, data, 0644)
}

func fixFileOwnership(configPath string) error {
	if os.Geteuid() != 0 {
		return nil
	}
	
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser == "" {
		return nil
	}
	
	uid, gid, err := getSudoUserIDs(sudoUser)
	if err != nil {
		return err
	}
	
	return applyOwnership(configPath, uid, gid)
}

func getSudoUserIDs(sudoUser string) (int, int, error) {
	realUser, err := user.Lookup(sudoUser)
	if err != nil {
		return 0, 0, err
	}
	
	uid, err := strconv.Atoi(realUser.Uid)
	if err != nil {
		return 0, 0, err
	}
	
	gid, err := strconv.Atoi(realUser.Gid)
	if err != nil {
		return 0, 0, err
	}
	
	return uid, gid, nil
}

func applyOwnership(configPath string, uid, gid int) error {
	if err := os.Chown(configPath, uid, gid); err != nil {
		return err
	}
	
	configDir := filepath.Dir(configPath)
	return os.Chown(configDir, uid, gid)
}

func (c *Config) AddInstalledPHP(version, installPath string) {
	if c.InstalledPHP == nil {
		c.InstalledPHP = make(map[string]PHPInfo)
	}
	
	c.InstalledPHP[version] = PHPInfo{
		Version:     version,
		InstallPath: installPath,
		InstallDate: time.Now(),
		IsCLI:       c.CurrentCLI == version,
		Extensions:  getDefaultExtensions(),
	}
}

func (c *Config) AddInstalledPHPWithExtensions(version, installPath string, extensions []string) {
	if c.InstalledPHP == nil {
		c.InstalledPHP = make(map[string]PHPInfo)
	}
	
	c.InstalledPHP[version] = PHPInfo{
		Version:     version,
		InstallPath: installPath,
		InstallDate: time.Now(),
		IsCLI:       c.CurrentCLI == version,
		Extensions:  extensions,
	}
}

func (c *Config) UpdatePHPExtensions(version string, extensions []string) bool {
	if info, exists := c.InstalledPHP[version]; exists {
		info.Extensions = extensions
		c.InstalledPHP[version] = info
		return true
	}
	return false
}

func (c *Config) GetPHPExtensions(version string) ([]string, bool) {
	if info, exists := c.InstalledPHP[version]; exists {
		return info.Extensions, true
	}
	return nil, false
}

func getDefaultExtensions() []string {
	return []string{
		"mbstring",
		"bcmath", 
		"opcache",
		"curl",
		"openssl",
		"zip",
		"sockets",
		"mysqli",
		"pdo-mysql",
		"gd",
		"jpeg",
		"freetype",
	}
}

// ConfigSnapshot represents a backup of configuration state
type ConfigSnapshot struct {
	Version    string
	Extensions []string
}

// CreateSnapshot creates a backup of the current extensions for a version
func (c *Config) CreateSnapshot(version string) *ConfigSnapshot {
	if info, exists := c.InstalledPHP[version]; exists {
		// Make a deep copy of the extensions slice
		extensionsCopy := make([]string, len(info.Extensions))
		copy(extensionsCopy, info.Extensions)
		
		return &ConfigSnapshot{
			Version:    version,
			Extensions: extensionsCopy,
		}
	}
	
	return &ConfigSnapshot{
		Version:    version,
		Extensions: []string{},
	}
}

// RestoreSnapshot restores configuration from a snapshot
func (c *Config) RestoreSnapshot(snapshot *ConfigSnapshot) {
	if info, exists := c.InstalledPHP[snapshot.Version]; exists {
		// Restore the extensions from the snapshot
		info.Extensions = make([]string, len(snapshot.Extensions))
		copy(info.Extensions, snapshot.Extensions)
		c.InstalledPHP[snapshot.Version] = info
	}
}

func (c *Config) RemoveInstalledPHP(version string) {
	if c.InstalledPHP != nil {
		delete(c.InstalledPHP, version)
	}
	
	if c.CurrentCLI == version {
		c.CurrentCLI = ""
	}
}

func (c *Config) SetCurrentCLI(version string) {
	if _, exists := c.InstalledPHP[version]; exists {
		c.CurrentCLI = version
		
		for v, info := range c.InstalledPHP {
			info.IsCLI = (v == version)
			c.InstalledPHP[v] = info
		}
	}
}

func getDefaultConfig() *Config {
	return &Config{
		YerdHome:       "/opt/yerd",
		InstalledPHP:   make(map[string]PHPInfo),
		CurrentCLI:     "",
		PathConfigured: false,
	}
}