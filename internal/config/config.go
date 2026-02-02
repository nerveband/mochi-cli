package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	configDir  = ".mochi-cli"
	configFile = "config.json"
)

// Profile represents a Mochi API profile
type Profile struct {
	APIKey string `json:"api_key"`
}

// Config represents the CLI configuration
type Config struct {
	ActiveProfile string             `json:"active_profile"`
	Profiles      map[string]Profile `json:"profiles"`
	DefaultFormat string             `json:"default_format,omitempty"`
	mu            sync.RWMutex
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig returns the singleton config instance
func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		instance, err = loadConfig()
	})
	return instance, err
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", configDir, configFile)
	}
	return filepath.Join(home, configDir, configFile)
}

// loadConfig loads configuration from disk
func loadConfig() (*Config, error) {
	path := getConfigPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config
		return &Config{
			Profiles:      make(map[string]Profile),
			DefaultFormat: "json",
		}, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.saveUnlocked()
}

// saveUnlocked saves without acquiring lock (caller must hold lock)
func (c *Config) saveUnlocked() error {
	path := getConfigPath()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Set secure permissions
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AddProfile adds a new profile
func (c *Config) AddProfile(name, apiKey string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Profiles[name] = Profile{
		APIKey: apiKey,
	}

	// If this is the first profile, set it as active
	if c.ActiveProfile == "" {
		c.ActiveProfile = name
	}

	return c.saveUnlocked()
}

// RemoveProfile removes a profile
func (c *Config) RemoveProfile(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Profiles[name]; !ok {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	delete(c.Profiles, name)

	// If we removed the active profile, clear it
	if c.ActiveProfile == name {
		c.ActiveProfile = ""
		// Set first available profile as active
		for n := range c.Profiles {
			c.ActiveProfile = n
			break
		}
	}

	return c.saveUnlocked()
}

// UseProfile sets the active profile
func (c *Config) UseProfile(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Profiles[name]; !ok {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	c.ActiveProfile = name
	return c.saveUnlocked()
}

// GetActiveProfile returns the active profile
func (c *Config) GetActiveProfile() (string, Profile, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.ActiveProfile == "" {
		return "", Profile{}, fmt.Errorf("no active profile set. Run 'mochi config add <name> <api-key>' to create a profile")
	}

	profile, ok := c.Profiles[c.ActiveProfile]
	if !ok {
		return "", Profile{}, fmt.Errorf("active profile '%s' not found", c.ActiveProfile)
	}

	return c.ActiveProfile, profile, nil
}

// GetProfile returns a specific profile
func (c *Config) GetProfile(name string) (Profile, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	profile, ok := c.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile '%s' not found", name)
	}

	return profile, nil
}

// ListProfiles returns all profile names
func (c *Config) ListProfiles() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.Profiles))
	for name := range c.Profiles {
		names = append(names, name)
	}

	return names
}

// Reset removes all configuration
func (c *Config) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Profiles = make(map[string]Profile)
	c.ActiveProfile = ""

	return c.saveUnlocked()
}

// GetAPIKey returns the API key, checking environment variable first
func GetAPIKey(profile Profile) string {
	// Check environment variable first
	if key := os.Getenv("MOCHI_API_KEY"); key != "" {
		return key
	}

	// Fall back to profile
	return profile.APIKey
}
