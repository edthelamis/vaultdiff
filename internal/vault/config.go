package vault

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvironmentConfig holds connection details for a single Vault environment.
type EnvironmentConfig struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	Mount   string `yaml:"mount"`
}

// Config holds the top-level vaultdiff configuration.
type Config struct {
	Environments []EnvironmentConfig `yaml:"environments"`
}

// LoadConfig reads and parses a YAML config file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// GetEnvironment returns the config for the named environment.
func (c *Config) GetEnvironment(name string) (*EnvironmentConfig, error) {
	for i := range c.Environments {
		if c.Environments[i].Name == name {
			return &c.Environments[i], nil
		}
	}
	return nil, fmt.Errorf("environment %q not found in config", name)
}

func (c *Config) validate() error {
	if len(c.Environments) == 0 {
		return fmt.Errorf("at least one environment must be defined")
	}
	seen := make(map[string]bool)
	for _, env := range c.Environments {
		if env.Name == "" {
			return fmt.Errorf("environment name must not be empty")
		}
		if seen[env.Name] {
			return fmt.Errorf("duplicate environment name %q", env.Name)
		}
		seen[env.Name] = true
		if env.Mount == "" {
			return fmt.Errorf("environment %q: mount must not be empty", env.Name)
		}
	}
	return nil
}
