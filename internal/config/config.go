package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/LumoSolutions/yerd/pkg/constants"
)

type Config struct {
	YerdHome       string             `json:"yerd_home"`
	InstalledPHP   map[string]PHPInfo `json:"installed_php"`
	CurrentCLI     string             `json:"current_cli"`
	PathConfigured bool               `json:"path_configured"`
}

type PHPInfo struct {
	Version     string    `json:"version"`
	InstallPath string    `json:"install_path"`
	InstallDate time.Time `json:"install_date"`
	IsCLI       bool      `json:"is_cli"`
	Extensions  []string  `json:"extensions"`
}

// GetConfigPath returns the path to the YERD configuration file, handling sudo user context.
// Creates config directory if it doesn't exist. Returns full path to config.json or error.
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

// LoadConfig loads YERD configuration from file or returns default config if file doesn't exist.
// Returns Config pointer or error if file read/parse fails.
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

// Save writes the configuration to file with proper JSON formatting and fixes ownership.
// Returns error if write fails or ownership fix fails.
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

// writeConfigFile marshals config to JSON and writes to specified path.
// configPath: Path to write config file. Returns error if marshal or write fails.
func (c *Config) writeConfigFile(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// fixFileOwnership changes config file ownership to sudo user when running as root.
// configPath: Path to config file. Returns error if ownership change fails.
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

// getSudoUserIDs retrieves UID and GID for the sudo user.
// sudoUser: Username of sudo user. Returns uid, gid integers or error if lookup fails.
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

// applyOwnership changes ownership of config file and directory to specified UID/GID.
// configPath: Path to config file, uid: User ID, gid: Group ID. Returns error if chown fails.
func applyOwnership(configPath string, uid, gid int) error {
	if err := os.Chown(configPath, uid, gid); err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	return os.Chown(configDir, uid, gid)
}

// AddInstalledPHP adds a new PHP installation to config with default extensions.
// version: PHP version string, installPath: Installation directory path.
func (c *Config) AddInstalledPHP(version, installPath string) {
	if c.InstalledPHP == nil {
		c.InstalledPHP = make(map[string]PHPInfo)
	}

	c.InstalledPHP[version] = PHPInfo{
		Version:     version,
		InstallPath: installPath,
		InstallDate: time.Now(),
		IsCLI:       c.CurrentCLI == version,
		Extensions:  constants.DefaultPHPExtensions,
	}
}

// AddInstalledPHPWithExtensions adds PHP installation with specified extensions.
// version: PHP version, installPath: Installation path, extensions: Extension list.
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

// UpdatePHPExtensions updates the extension list for a specific PHP version.
// version: PHP version string, extensions: New extension list. Returns true if version exists.
func (c *Config) UpdatePHPExtensions(version string, extensions []string) bool {
	if info, exists := c.InstalledPHP[version]; exists {
		info.Extensions = extensions
		c.InstalledPHP[version] = info
		return true
	}
	return false
}

// GetPHPExtensions retrieves extension list for a specific PHP version.
// version: PHP version string. Returns extension slice and existence boolean.
func (c *Config) GetPHPExtensions(version string) ([]string, bool) {
	if info, exists := c.InstalledPHP[version]; exists {
		return info.Extensions, true
	}
	return nil, false
}

type ConfigSnapshot struct {
	Version    string
	Extensions []string
}

// CreateSnapshot creates a backup of current configuration for rollback purposes.
// version: PHP version to snapshot. Returns ConfigSnapshot with extension list copy.
func (c *Config) CreateSnapshot(version string) *ConfigSnapshot {
	if info, exists := c.InstalledPHP[version]; exists {
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

// RestoreSnapshot restores configuration from a previous snapshot.
// snapshot: ConfigSnapshot containing version and extension data to restore.
func (c *Config) RestoreSnapshot(snapshot *ConfigSnapshot) {
	if info, exists := c.InstalledPHP[snapshot.Version]; exists {
		info.Extensions = make([]string, len(snapshot.Extensions))
		copy(info.Extensions, snapshot.Extensions)
		c.InstalledPHP[snapshot.Version] = info
	}
}

// RemoveInstalledPHP removes a PHP version from config and clears CLI if it was current.
// version: PHP version string to remove from installed versions.
func (c *Config) RemoveInstalledPHP(version string) {
	if c.InstalledPHP != nil {
		delete(c.InstalledPHP, version)
	}

	if c.CurrentCLI == version {
		c.CurrentCLI = ""
	}
}

// SetCurrentCLI sets the specified PHP version as the current CLI and updates IsCLI flags.
// version: PHP version string to set as current CLI version.
func (c *Config) SetCurrentCLI(version string) {
	if _, exists := c.InstalledPHP[version]; exists {
		c.CurrentCLI = version

		for v, info := range c.InstalledPHP {
			info.IsCLI = (v == version)
			c.InstalledPHP[v] = info
		}
	}
}

// getDefaultConfig returns a new Config with default values and empty PHP installations.
func getDefaultConfig() *Config {
	return &Config{
		YerdHome:       "/opt/yerd",
		InstalledPHP:   make(map[string]PHPInfo),
		CurrentCLI:     "",
		PathConfigured: false,
	}
}
