package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved reference state of secrets for an environment.
type Baseline struct {
	Environment string            `json:"environment"`
	Path        string            `json:"path"`
	CreatedAt   time.Time         `json:"created_at"`
	Data        map[string]string `json:"data"`
}

// SaveBaseline writes a baseline to disk at the given file path.
func SaveBaseline(b *Baseline, filePath string) error {
	if b == nil {
		return fmt.Errorf("baseline is nil")
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create baseline file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("failed to encode baseline: %w", err)
	}
	return nil
}

// LoadBaseline reads a baseline from disk.
func LoadBaseline(filePath string) (*Baseline, error) {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline file not found: %s", filePath)
		}
		return nil, fmt.Errorf("failed to open baseline file: %w", err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("failed to decode baseline: %w", err)
	}
	return &b, nil
}

// NewBaseline creates a new Baseline from a snapshot's data.
func NewBaseline(environment, path string, data map[string]string) *Baseline {
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return &Baseline{
		Environment: environment,
		Path:        path,
		CreatedAt:   time.Now().UTC(),
		Data:        copy,
	}
}
