package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Source      string   `yaml:"source"`
	Destination string   `yaml:"destination"`
	Workers     int      `yaml:"workers"`
	Overwrite   bool     `yaml:"overwrite"`
	Extensions  []string `yaml:"extensions"`
	MaxRetries  int      `yaml:"max_retries"`
	DryRun      bool     `yaml:"dry_run"`
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		Source:      "",
		Destination: "",
		Workers:     10,
		Overwrite:   false,
		Extensions:  []string{},
		MaxRetries:  3,
		DryRun:      false,
	}
}

// LoadFromFile loads configuration from a YAML file
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

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Source == "" {
		return fmt.Errorf("source path is required")
	}
	if c.Destination == "" {
		return fmt.Errorf("destination path is required")
	}
	if c.Workers < 1 {
		c.Workers = 1
	}
	if c.Workers > 50 {
		c.Workers = 50
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}
	return nil
}

// HasExtensionFilter checks if extension filtering is enabled
func (c *Config) HasExtensionFilter() bool {
	return len(c.Extensions) > 0
}

// IsExtensionAllowed checks if a file extension is in the allowed list
func (c *Config) IsExtensionAllowed(ext string) bool {
	if !c.HasExtensionFilter() {
		return true
	}
	for _, allowed := range c.Extensions {
		if allowed == ext {
			return true
		}
	}
	return false
}
