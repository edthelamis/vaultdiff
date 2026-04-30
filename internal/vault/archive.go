package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ArchiveEntry represents a single archived snapshot of a secret path.
type ArchiveEntry struct {
	Path      string                 `json:"path"`
	Data      map[string]interface{} `json:"data"`
	ArchivedAt time.Time             `json:"archived_at"`
	Label     string                 `json:"label,omitempty"`
}

// Archive holds a collection of archived secret snapshots.
type Archive struct {
	Entries []ArchiveEntry `json:"entries"`
}

// NewArchive creates an empty Archive.
func NewArchive() *Archive {
	return &Archive{Entries: []ArchiveEntry{}}
}

// Add appends a new entry to the archive for the given path and data.
func (a *Archive) Add(path string, data map[string]interface{}, label string) {
	copy := make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	a.Entries = append(a.Entries, ArchiveEntry{
		Path:       path,
		Data:       copy,
		ArchivedAt: time.Now().UTC(),
		Label:      label,
	})
}

// FilterByPath returns all entries matching the given path.
func (a *Archive) FilterByPath(path string) []ArchiveEntry {
	var result []ArchiveEntry
	for _, e := range a.Entries {
		if e.Path == path {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a human-readable summary of the archive.
func (a *Archive) Summary() string {
	paths := make(map[string]int)
	for _, e := range a.Entries {
		paths[e.Path]++
	}
	keys := make([]string, 0, len(paths))
	for k := range paths {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	summary := fmt.Sprintf("Archive: %d entries across %d path(s)\n", len(a.Entries), len(paths))
	for _, k := range keys {
		summary += fmt.Sprintf("  %s: %d version(s)\n", k, paths[k])
	}
	return summary
}

// SaveArchive writes the archive to a JSON file at the given path.
func SaveArchive(a *Archive, filePath string) error {
	if a == nil {
		return fmt.Errorf("archive is nil")
	}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal archive: %w", err)
	}
	return os.WriteFile(filePath, data, 0600)
}

// LoadArchive reads an archive from a JSON file.
func LoadArchive(filePath string) (*Archive, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read archive file: %w", err)
	}
	var a Archive
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("unmarshal archive: %w", err)
	}
	return &a, nil
}
