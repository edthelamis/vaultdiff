package vault

import (
	"encoding/json"
	"fmt"
	"os"
)

// ArchiveStoreOptions configures archive persistence.
type ArchiveStoreOptions struct {
	// FilePath is the path to the archive JSON file.
	FilePath string
}

// SaveArchiveStore writes an Archive to disk using the provided options.
func SaveArchiveStore(a *Archive, opts ArchiveStoreOptions) error {
	if opts.FilePath == "" {
		return fmt.Errorf("archive store file path must not be empty")
	}
	return SaveArchive(a, opts.FilePath)
}

// LoadArchiveStore reads an Archive from disk using the provided options.
// If the file does not exist, an empty archive is returned.
func LoadArchiveStore(opts ArchiveStoreOptions) (*Archive, error) {
	if opts.FilePath == "" {
		return nil, fmt.Errorf("archive store file path must not be empty")
	}
	if _, err := os.Stat(opts.FilePath); os.IsNotExist(err) {
		return NewArchive(), nil
	}
	data, err := os.ReadFile(opts.FilePath)
	if err != nil {
		return nil, fmt.Errorf("read archive store: %w", err)
	}
	var a Archive
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("unmarshal archive store: %w", err)
	}
	return &a, nil
}
