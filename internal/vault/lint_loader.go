package vault

import (
	"encoding/json"
	"fmt"
	"os"
)

// LintConfig holds a named set of lint rules loaded from a file.
type LintConfig struct {
	Name  string     `json:"name"`
	Rules []LintRule `json:"rules"`
}

// LoadLintConfig reads and parses a JSON lint config from the given path.
func LoadLintConfig(path string) (*LintConfig, error) {
	if path == "" {
		return nil, fmt.Errorf("lint config path must not be empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read lint config: %w", err)
	}
	var cfg LintConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse lint config: %w", err)
	}
	if cfg.Name == "" {
		return nil, fmt.Errorf("lint config missing required field: name")
	}
	if len(cfg.Rules) == 0 {
		return nil, fmt.Errorf("lint config %q has no rules defined", cfg.Name)
	}
	for i, r := range cfg.Rules {
		if r.Name == "" {
			return nil, fmt.Errorf("lint rule at index %d missing name", i)
		}
		if r.Pattern == "" {
			return nil, fmt.Errorf("lint rule %q missing pattern", r.Name)
		}
		if r.Target != "key" && r.Target != "value" {
			return nil, fmt.Errorf("lint rule %q has invalid target %q (must be 'key' or 'value')", r.Name, r.Target)
		}
	}
	return &cfg, nil
}
