package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveSignatureRecord persists a SignatureRecord to a JSON file.
func SaveSignatureRecord(path string, record *SignatureRecord) error {
	if record == nil {
		return fmt.Errorf("save signature: record is nil")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("save signature: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("save signature: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(record); err != nil {
		return fmt.Errorf("save signature: encode: %w", err)
	}
	return nil
}

// LoadSignatureRecord reads a SignatureRecord from a JSON file.
func LoadSignatureRecord(path string) (*SignatureRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("load signature: file not found: %s", path)
		}
		return nil, fmt.Errorf("load signature: open: %w", err)
	}
	defer f.Close()
	var record SignatureRecord
	if err := json.NewDecoder(f).Decode(&record); err != nil {
		return nil, fmt.Errorf("load signature: decode: %w", err)
	}
	return &record, nil
}
