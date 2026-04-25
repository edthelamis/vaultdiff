package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveTagIndex persists a TagIndex to a JSON file at the given path.
func SaveTagIndex(idx TagIndex, path string) error {
	if idx == nil {
		return fmt.Errorf("tag index is nil")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating tag index file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(idx); err != nil {
		return fmt.Errorf("encoding tag index: %w", err)
	}
	return nil
}

// LoadTagIndex reads a TagIndex from a JSON file at the given path.
func LoadTagIndex(path string) (TagIndex, error) {
	if path == "" {
		return nil, fmt.Errorf("tag index path is empty")
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("tag index file not found: %s", path)
		}
		return nil, fmt.Errorf("opening tag index file: %w", err)
	}
	defer f.Close()
	var idx TagIndex
	if err := json.NewDecoder(f).Decode(&idx); err != nil {
		return nil, fmt.Errorf("decoding tag index: %w", err)
	}
	return idx, nil
}
