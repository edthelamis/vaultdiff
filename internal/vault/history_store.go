package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// SaveHistory serialises a History to a JSON file at the given path.
func SaveHistory(h *History, path string) error {
	if h == nil {
		return errors.New("history is nil")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating history directory: %w", err)
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling history: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing history file: %w", err)
	}
	return nil
}

// LoadHistory deserialises a History from a JSON file at the given path.
func LoadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("history file not found: %s", path)
		}
		return nil, fmt.Errorf("reading history file: %w", err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("parsing history file: %w", err)
	}
	return &h, nil
}
