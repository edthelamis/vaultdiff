package vault

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveAccessLog writes an AccessLog to a JSON file at the given path.
func SaveAccessLog(log *AccessLog, path string) error {
	if log == nil {
		return fmt.Errorf("access log is nil")
	}
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal access log: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write access log: %w", err)
	}
	return nil
}

// LoadAccessLog reads an AccessLog from a JSON file at the given path.
func LoadAccessLog(path string) (*AccessLog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("access log file not found: %s", path)
		}
		return nil, fmt.Errorf("read access log: %w", err)
	}
	var log AccessLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse access log: %w", err)
	}
	return &log, nil
}
