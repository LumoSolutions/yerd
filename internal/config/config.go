package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lumosolutions/yerd/internal/constants"
	"github.com/lumosolutions/yerd/internal/utils"
)

var (
	instance  *Config
	once      sync.Once
	initErr   error
	pathMutex sync.RWMutex
)

// Config represents the configuration manager
//
// Path notation supports both dot-separated paths and bracket literals:
//   - Standard: "installed.php.version" creates nested objects
//   - Brackets: "installed.[php.8.3]" treats "php.8.3" as a single key
//   - Mixed: "config.[db.primary].host" combines both approaches
//
// The default config path is "config.json" in the current directory.
// Use SetConfigPath() before any other operations to change it.
type Config struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	filePath string
}

// getInstance returns the singleton config instance, initializing it if needed
func getInstance() (*Config, error) {
	once.Do(func() {
		pathMutex.RLock()
		configDir, _ := utils.GetUserConfigDir()
		path := filepath.Join(configDir, constants.YerdConfigName)
		pathMutex.RUnlock()

		instance = &Config{
			filePath: path,
			data:     make(map[string]interface{}),
		}

		if err := instance.Load(); err != nil {
			if !os.IsNotExist(err) {
				initErr = fmt.Errorf("failed to load config: %w", err)
			}
		}
	})

	return instance, initErr
}

// GetConfig retrieves a configuration value at the specified path
// Path can use dots as separators or brackets to preserve dots in keys:
func getConfig(name string) (interface{}, error) {
	cfg, err := getInstance()
	if err != nil {
		return nil, err
	}

	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	parts := parsePath(name)
	return getValueAtPath(cfg.data, parts)
}

// GetString retrieves a string value at the specified path
func GetString(name string) (string, error) {
	val, err := getConfig(name)
	if err != nil {
		return "", err
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value at path %s is not a string", name)
	}

	return str, nil
}

// GetObject retrieves an object (map) at the specified path
func GetObject(name string) (map[string]interface{}, error) {
	val, err := getConfig(name)
	if err != nil {
		return nil, err
	}

	obj, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value at path %s is not an object", name)
	}

	return obj, nil
}

// SetObjectConfig sets an object at the specified path
// Automatically creates parent nodes if they don't exist
// Use brackets to preserve dots: "installed.[php.8.3]" creates {"installed": {"php.8.3": data}}
func SetObjectConfig(name string, data interface{}) error {
	cfg, err := getInstance()
	if err != nil {
		utils.LogError(err, "config")
		return err
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Convert data to map if it's a struct or other type
	var objData map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		objData = v
	default:
		// Try to convert via JSON marshaling/unmarshaling
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			utils.LogError(err, "config")
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &objData); err != nil {
			utils.LogError(err, "config")
			return fmt.Errorf("failed to unmarshal data to map: %w", err)
		}
	}

	parts := parsePath(name)
	setValueAtPath(cfg.data, parts, objData)

	// Auto-save after modification
	return cfg.save()
}

// SetStringData sets a string value at the specified path
// Automatically creates parent nodes if they don't exist
// Use brackets to preserve dots: "installed.[php.8.3]" creates {"installed": {"php.8.3": data}}
func SetStringData(name string, data string) error {
	cfg, err := getInstance()
	if err != nil {
		return err
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	parts := parsePath(name)
	setValueAtPath(cfg.data, parts, data)

	// Auto-save after modification
	return cfg.save()
}

// SetStruct sets a struct value at the specified path
// The struct will be converted to a map structure in the config
func SetStruct(name string, value interface{}) error {
	return SetObjectConfig(name, value)
}

// GetStruct retrieves configuration at the path and unmarshals it into the provided struct
// The dest parameter must be a pointer to a struct
func GetStruct(name string, dest interface{}) error {
	val, err := getConfig(name)
	if err != nil {
		utils.LogInfo("getstruct", "Not found")
		utils.LogError(err, "getstruct")
		return err
	}

	// Convert via JSON marshaling/unmarshaling
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		utils.LogError(err, "getstruct")
		return fmt.Errorf("failed to marshal config value: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, dest); err != nil {
		utils.LogError(err, "getstruct")
		return fmt.Errorf("failed to unmarshal into struct: %w", err)
	}

	return nil
}

// Delete removes a configuration value at the specified path
func Delete(name string) error {
	cfg, err := getInstance()
	if err != nil {
		return err
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	parts := parsePath(name)
	if err := deleteValueAtPath(cfg.data, parts); err != nil {
		return err
	}

	// Auto-save after modification
	return cfg.save()
}

// Exists checks if a configuration exists at the specified path
func Exists(name string) bool {
	_, err := getConfig(name)
	return err == nil
}

// GetAll returns the entire configuration data
func GetAll() (map[string]interface{}, error) {
	cfg, err := getInstance()
	if err != nil {
		return nil, err
	}

	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	// Return a deep copy to prevent external modifications
	return deepCopy(cfg.data).(map[string]interface{}), nil
}

// Clear removes all configuration data
func Clear() error {
	cfg, err := getInstance()
	if err != nil {
		return err
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.data = make(map[string]interface{})

	// Auto-save after modification
	return cfg.save()
}

// Save explicitly saves the configuration to file (usually auto-saved)
func Save() error {
	cfg, err := getInstance()
	if err != nil {
		return err
	}

	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	return cfg.save()
}

// Load explicitly reloads the configuration from file
func Load() error {
	cfg, err := getInstance()
	if err != nil {
		return err
	}

	return cfg.Load()
}

// Internal methods for Config struct

// Load reads the configuration from the file
func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		c.data = make(map[string]interface{})
		return nil
	}

	if err := json.Unmarshal(data, &c.data); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// save writes the configuration to the file (internal, no lock)
func (c *Config) save() error {
	data, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	userCtx, err := utils.GetRealUser()
	if err == nil {
		utils.Chown(c.filePath, userCtx.UID, userCtx.GID)
	}

	return nil
}

// parsePath splits the dot-notation path into parts, respecting bracket literals
func parsePath(path string) []string {
	parser := &pathParser{
		parts:   []string{},
		current: strings.Builder{},
	}

	for i := 0; i < len(path); i++ {
		parser.processChar(path[i])
	}

	parser.finalize()
	return parser.parts
}

// pathParser encapsulates the parsing state and logic
type pathParser struct {
	parts     []string
	current   strings.Builder
	inBracket bool
}

// processChar handles a single character based on current state
func (p *pathParser) processChar(ch byte) {
	switch ch {
	case '[':
		p.handleOpenBracket(ch)
	case ']':
		p.handleCloseBracket(ch)
	case '.':
		p.handleDot(ch)
	default:
		p.current.WriteByte(ch)
	}
}

// handleOpenBracket processes '[' character
func (p *pathParser) handleOpenBracket(ch byte) {
	if p.inBracket {
		// Nested bracket, treat as literal
		p.current.WriteByte(ch)
		return
	}

	// Start of bracket literal
	p.inBracket = true
	p.saveCurrentPart()
}

// handleCloseBracket processes ']' character
func (p *pathParser) handleCloseBracket(ch byte) {
	if !p.inBracket {
		// Stray closing bracket, treat as literal
		p.current.WriteByte(ch)
		return
	}

	// End of bracket literal
	p.inBracket = false
	p.saveCurrentPart()
}

// handleDot processes '.' character
func (p *pathParser) handleDot(ch byte) {
	if p.inBracket {
		// Inside brackets, dots are literal
		p.current.WriteByte(ch)
		return
	}

	// Outside brackets, dot is a separator
	p.saveCurrentPart()
}

// saveCurrentPart adds current content to parts if not empty and resets builder
func (p *pathParser) saveCurrentPart() {
	if p.current.Len() > 0 {
		p.parts = append(p.parts, p.current.String())
		p.current.Reset()
	}
}

// finalize adds any remaining content to parts
func (p *pathParser) finalize() {
	p.saveCurrentPart()
}

// getValueAtPath retrieves a value from nested maps using a path
func getValueAtPath(data map[string]interface{}, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - return the value
			val, exists := current[part]
			if !exists {
				return nil, fmt.Errorf("path not found: %s", strings.Join(parts[:i+1], "."))
			}
			return val, nil
		}

		// Navigate deeper
		next, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("path not found: %s", strings.Join(parts[:i+1], "."))
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path %s is not an object", strings.Join(parts[:i+1], "."))
		}

		current = nextMap
	}

	return nil, fmt.Errorf("unexpected end of path")
}

// setValueAtPath sets a value in nested maps using a path, creating parents as needed
func setValueAtPath(data map[string]interface{}, parts []string, value interface{}) {
	if len(parts) == 0 {
		return
	}

	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - set the value
			current[part] = value
			return
		}

		// Navigate or create deeper structure
		next, exists := current[part]
		if !exists {
			// Create new map
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		} else {
			// Check if it's a map
			nextMap, ok := next.(map[string]interface{})
			if !ok {
				// Overwrite non-map value with a new map
				newMap := make(map[string]interface{})
				current[part] = newMap
				current = newMap
			} else {
				current = nextMap
			}
		}
	}
}

// deleteValueAtPath removes a value from nested maps using a path
func deleteValueAtPath(data map[string]interface{}, parts []string) error {
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - delete the value
			delete(current, part)
			return nil
		}

		// Navigate deeper
		next, exists := current[part]
		if !exists {
			return fmt.Errorf("path not found: %s", strings.Join(parts[:i+1], "."))
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("path %s is not an object", strings.Join(parts[:i+1], "."))
		}

		current = nextMap
	}

	return fmt.Errorf("unexpected end of path")
}

// deepCopy creates a deep copy of the data structure
func deepCopy(src interface{}) interface{} {
	switch v := src.(type) {
	case map[string]interface{}:
		dst := make(map[string]interface{})
		for key, val := range v {
			dst[key] = deepCopy(val)
		}
		return dst
	case []interface{}:
		dst := make([]interface{}, len(v))
		for i, val := range v {
			dst[i] = deepCopy(val)
		}
		return dst
	default:
		return v
	}
}
