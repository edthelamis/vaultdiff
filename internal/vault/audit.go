package vault

import (
	"fmt"
	"time"
)

// AuditEntry represents a single recorded diff audit event.
type AuditEntry struct {
	Timestamp   time.Time
	Environment string
	Path        string
	VersionA    int
	VersionB    int
	Changes     []DiffEntry
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new audit entry to the log.
func (a *AuditLog) Record(env, path string, versionA, versionB int, changes []DiffEntry) {
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Path:        path,
		VersionA:    versionA,
		VersionB:    versionB,
		Changes:     changes,
	})
}

// Summary returns a human-readable summary of all recorded audit entries.
func (a *AuditLog) Summary() string {
	if len(a.Entries) == 0 {
		return "No audit entries recorded."
	}
	var out string
	for _, e := range a.Entries {
		out += fmt.Sprintf("[%s] env=%s path=%s v%d..v%d changes=%d\n",
			e.Timestamp.Format(time.RFC3339),
			e.Environment,
			e.Path,
			e.VersionA,
			e.VersionB,
			len(e.Changes),
		)
	}
	return out
}

// HasChanges returns true if any recorded entry contains at least one change.
func (a *AuditLog) HasChanges() bool {
	for _, e := range a.Entries {
		if len(e.Changes) > 0 {
			return true
		}
	}
	return false
}
