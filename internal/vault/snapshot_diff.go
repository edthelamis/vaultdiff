package vault

import (
	"fmt"
)

// SnapshotDiffResult holds the diff result between two snapshots along with metadata.
type SnapshotDiffResult struct {
	FromSnapshot *Snapshot
	ToSnapshot   *Snapshot
	Changes      []DiffEntry
}

// DiffSnapshots computes the diff between two Snapshot values.
// It returns a SnapshotDiffResult containing all changed entries.
func DiffSnapshots(from, to *Snapshot) (*SnapshotDiffResult, error) {
	if from == nil || to == nil {
		return nil, fmt.Errorf("snapshot_diff: both snapshots must be non-nil")
	}
	changes := DiffSecrets(from.Data, to.Data)
	return &SnapshotDiffResult{
		FromSnapshot: from,
		ToSnapshot:   to,
		Changes:      changes,
	}, nil
}

// HasChanges returns true if the diff result contains any non-unchanged entries.
func (r *SnapshotDiffResult) HasChanges() bool {
	for _, c := range r.Changes {
		if c.ChangeType != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable one-line summary of the diff.
func (r *SnapshotDiffResult) Summary() string {
	var added, removed, changed int
	for _, c := range r.Changes {
		switch c.ChangeType {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf(
		"env=%s path=%s v%d→v%d: +%d -%d ~%d",
		r.FromSnapshot.Environment,
		r.FromSnapshot.Path,
		r.FromSnapshot.Version,
		r.ToSnapshot.Version,
		added, removed, changed,
	)
}
