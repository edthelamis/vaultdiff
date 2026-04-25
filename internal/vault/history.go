package vault

import (
	"fmt"
	"sort"
	"time"
)

// HistoryEntry represents a single recorded diff event between two secret versions.
type HistoryEntry struct {
	Timestamp   time.Time        `json:"timestamp"`
	Environment string           `json:"environment"`
	Path        string           `json:"path"`
	FromVersion int              `json:"from_version"`
	ToVersion   int              `json:"to_version"`
	Changes     []DiffEntry      `json:"changes"`
}

// History holds an ordered list of diff history entries.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// NewHistory creates an empty History.
func NewHistory() *History {
	return &History{Entries: []HistoryEntry{}}
}

// Record appends a new entry to the history.
func (h *History) Record(env, path string, from, to int, changes []DiffEntry) {
	h.Entries = append(h.Entries, HistoryEntry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Path:        path,
		FromVersion: from,
		ToVersion:   to,
		Changes:     changes,
	})
}

// Filter returns history entries matching the given environment and path.
// Pass empty strings to skip filtering on that field.
func (h *History) Filter(env, path string) []HistoryEntry {
	var result []HistoryEntry
	for _, e := range h.Entries {
		if (env == "" || e.Environment == env) && (path == "" || e.Path == path) {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a human-readable summary of the history.
func (h *History) Summary() string {
	if len(h.Entries) == 0 {
		return "No history recorded."
	}
	// Sort by timestamp ascending
	sorted := make([]HistoryEntry, len(h.Entries))
	copy(sorted, h.Entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})
	var out string
	for _, e := range sorted {
		out += fmt.Sprintf("[%s] %s %s v%d→v%d (%d changes)\n",
			e.Timestamp.Format(time.RFC3339),
			e.Environment,
			e.Path,
			e.FromVersion,
			e.ToVersion,
			len(e.Changes),
		)
	}
	return out
}
