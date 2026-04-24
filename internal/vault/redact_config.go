package vault

import (
	"encoding/json"
	"fmt"
	"os"
)

// RedactConfig is the on-disk representation of redaction settings.
type RedactConfig struct {
	KeyPatterns []string `json:"key_patterns"`
	Replacement string   `json:"replacement"`
}

// LoadRedactConfig reads a JSON redaction config file from the given path.
// If the file does not exist, DefaultRedactOptions is returned.
func LoadRedactConfig(path string) (RedactOptions, error) {
	if path == "" {
		return DefaultRedactOptions(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultRedactOptions(), nil
		}
		return RedactOptions{}, fmt.Errorf("reading redact config: %w", err)
	}
	var cfg RedactConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return RedactOptions{}, fmt.Errorf("parsing redact config: %w", err)
	}
	if len(cfg.KeyPatterns) == 0 {
		return RedactOptions{}, fmt.Errorf("redact config must include at least one key_pattern")
	}
	replacement := cfg.Replacement
	if replacement == "" {
		replacement = "[REDACTED]"
	}
	return RedactOptions{
		KeyPatterns: cfg.KeyPatterns,
		Replacement: replacement,
	}, nil
}
