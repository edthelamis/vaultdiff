package vault

import (
	"fmt"
	"sort"
	"time"
)

// AccessEntry records a single access event for a secret path.
type AccessEntry struct {
	Path      string    `json:"path"`
	Operation string    `json:"operation"` // read, write, delete
	Actor     string    `json:"actor"`
	Timestamp time.Time `json:"timestamp"`
}

// AccessLog holds all recorded access events.
type AccessLog struct {
	Entries []AccessEntry `json:"entries"`
}

// NewAccessLog creates an empty AccessLog.
func NewAccessLog() *AccessLog {
	return &AccessLog{Entries: []AccessEntry{}}
}

// Record appends a new access entry to the log.
func (a *AccessLog) Record(path, operation, actor string) {
	a.Entries = append(a.Entries, AccessEntry{
		Path:      path,
		Operation: operation,
		Actor:     actor,
		Timestamp: time.Now().UTC(),
	})
}

// FilterByActor returns entries matching the given actor.
func (a *AccessLog) FilterByActor(actor string) []AccessEntry {
	var out []AccessEntry
	for _, e := range a.Entries {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}

// FilterByPath returns entries matching the given path.
func (a *AccessLog) FilterByPath(path string) []AccessEntry {
	var out []AccessEntry
	for _, e := range a.Entries {
		if e.Path == path {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable summary of access counts per path.
func (a *AccessLog) Summary() string {
	if len(a.Entries) == 0 {
		return "no access events recorded"
	}
	counts := make(map[string]int)
	for _, e := range a.Entries {
		counts[e.Path]++
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := fmt.Sprintf("%d total access events:\n", len(a.Entries))
	for _, k := range keys {
		out += fmt.Sprintf("  %s: %d\n", k, counts[k])
	}
	return out
}
