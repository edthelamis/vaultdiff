package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveAnnotationIndex persists an AnnotationIndex to a JSON file.
func SaveAnnotationIndex(idx *AnnotationIndex, path string) error {
	if idx == nil {
		return fmt.Errorf("annotation index is nil")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling annotation index: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing annotation index: %w", err)
	}
	return nil
}

// LoadAnnotationIndex reads an AnnotationIndex from a JSON file.
func LoadAnnotationIndex(path string) (*AnnotationIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("annotation index file not found: %s", path)
		}
		return nil, fmt.Errorf("reading annotation index: %w", err)
	}
	var idx AnnotationIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parsing annotation index: %w", err)
	}
	if idx.Annotations == nil {
		idx.Annotations = make(map[string]Annotation)
	}
	return &idx, nil
}
