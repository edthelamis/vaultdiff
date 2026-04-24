package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures a point-in-time view of a secret's key-value pairs.
type Snapshot struct {
	Environment string            `json:"environment"`
	Path        string            `json:"path"`
	Version     int               `json:"version"`
	CapturedAt  time.Time         `json:"captured_at"`
	Data        map[string]string `json:"data"`
}

// TakeSnapshot records the current state of a secret at the given path.
func TakeSnapshot(env, path string, version int, data map[string]string) *Snapshot {
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return &Snapshot{
		Environment: env,
		Path:        path,
		Version:     version,
		CapturedAt:  time.Now().UTC(),
		Data:        copy,
	}
}

// SaveSnapshot writes a snapshot to a JSON file at the given filepath.
func SaveSnapshot(s *Snapshot, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from a JSON file at the given filepath.
func LoadSnapshot(filepath string) (*Snapshot, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
