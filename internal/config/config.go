package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Destination represents a single destination with its own settings.
// Each destination can have independent overwrite settings, allowing
// fine-grained control over how files are copied to different locations.
type Destination struct {
	ID        string `yaml:"id" json:"id"`
	Path      string `yaml:"path" json:"path"`
	Overwrite bool   `yaml:"overwrite" json:"overwrite"`
	Enabled   bool   `yaml:"enabled" json:"enabled"`
}

// CopyGroup represents a copy configuration with one source and multiple destinations.
// This enables the common use case of backing up/distributing files to multiple locations.
type CopyGroup struct {
	ID           string        `yaml:"id" json:"id"`
	Name         string        `yaml:"name" json:"name"`
	Source       string        `yaml:"source" json:"source"`
	Destinations []Destination `yaml:"destinations" json:"destinations"`
	Enabled      bool          `yaml:"enabled" json:"enabled"`
}

// Config represents the application configuration.
// It supports both legacy single source/destination mode and the new Copy Groups feature.
// JSON tags are added for Wails frontend binding.
type Config struct {
	// Legacy single source/destination (for backward compatibility with CLI mode)
	Source      string `yaml:"source" json:"source"`
	Destination string `yaml:"destination" json:"destination"`

	// Copy Groups - allows one source to copy to multiple destinations
	Groups []CopyGroup `yaml:"groups,omitempty" json:"groups"`

	// Global settings applied to all copy operations
	Workers    int      `yaml:"workers" json:"workers"`
	Overwrite  bool     `yaml:"overwrite" json:"overwrite"`
	Extensions []string `yaml:"extensions" json:"extensions"`
	MaxRetries int      `yaml:"max_retries" json:"maxRetries"`
	DryRun     bool     `yaml:"dry_run" json:"dryRun"`
}

// DefaultConfig returns a config with sensible default values.
// These defaults provide a good balance between performance and resource usage.
func DefaultConfig() *Config {
	return &Config{
		Source:      "",
		Destination: "",
		Groups:      []CopyGroup{},
		Workers:     10, // 10 concurrent workers is typically optimal for network file operations
		Overwrite:   false,
		Extensions:  []string{},
		MaxRetries:  3, // 3 retries with exponential backoff handles most transient failures
		DryRun:      false,
	}
}

// LoadFromFile loads configuration from a YAML file.
// Returns an error if the file cannot be read or parsed.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveToFile persists the configuration to a YAML file.
// This allows user preferences to survive application restarts.
func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Write config file with restricted permissions.
	// Using 0600 for security (only owner can read/write).
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid for copy operations.
// It also normalizes values to ensure they're within acceptable ranges.
func (c *Config) Validate() error {
	// In legacy mode, source and destination are required
	if len(c.Groups) == 0 {
		if c.Source == "" {
			return fmt.Errorf("source path is required")
		}
		if c.Destination == "" {
			return fmt.Errorf("destination path is required")
		}
	}

	// Clamp workers to a reasonable range.
	// Too few workers underutilizes resources; too many causes contention.
	if c.Workers < 1 {
		c.Workers = 1
	}
	if c.Workers > 50 {
		c.Workers = 50
	}

	// Negative retries don't make sense
	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}

	return nil
}

// HasExtensionFilter checks if extension filtering is enabled.
// When enabled, only files with matching extensions will be copied.
func (c *Config) HasExtensionFilter() bool {
	return len(c.Extensions) > 0
}

// IsExtensionAllowed checks if a file extension is in the allowed list.
// Returns true if no filter is set (all extensions allowed) or if the
// extension matches one in the allowed list.
func (c *Config) IsExtensionAllowed(ext string) bool {
	if !c.HasExtensionFilter() {
		return true
	}

	// Normalize the extension to lowercase for case-insensitive matching
	ext = strings.ToLower(ext)
	for _, allowed := range c.Extensions {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}
	return false
}

// GetEnabledGroups returns only the groups that are enabled.
// This is used when processing copy operations to skip disabled groups.
func (c *Config) GetEnabledGroups() []CopyGroup {
	var enabled []CopyGroup
	for _, g := range c.Groups {
		if g.Enabled {
			enabled = append(enabled, g)
		}
	}
	return enabled
}

// AddGroup adds a new copy group to the configuration.
// The group ID should be unique to allow proper identification.
func (c *Config) AddGroup(group CopyGroup) {
	c.Groups = append(c.Groups, group)
}

// RemoveGroup removes a group by its ID.
// Returns true if a group was removed, false if the ID was not found.
func (c *Config) RemoveGroup(groupID string) bool {
	for i, g := range c.Groups {
		if g.ID == groupID {
			c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
			return true
		}
	}
	return false
}

// FindGroup finds a group by its ID.
// Returns nil if no group with the given ID exists.
func (c *Config) FindGroup(groupID string) *CopyGroup {
	for i := range c.Groups {
		if c.Groups[i].ID == groupID {
			return &c.Groups[i]
		}
	}
	return nil
}
