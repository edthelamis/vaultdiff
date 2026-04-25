package vault

import "fmt"

// BaselineDiffResult holds the result of comparing current secrets against a baseline.
type BaselineDiffResult struct {
	Environment string
	Path        string
	Changes     []DiffEntry
}

// Summary returns a human-readable summary of baseline diff results.
func (r *BaselineDiffResult) Summary() string {
	if len(r.Changes) == 0 {
		return fmt.Sprintf("[%s] %s: no drift detected", r.Environment, r.Path)
	}
	added, removed, changed := 0, 0, 0
	for _, c := range r.Changes {
		switch c.Type {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeModified:
			changed++
		}
	}
	return fmt.Sprintf("[%s] %s: %d added, %d removed, %d changed",
		r.Environment, r.Path, added, removed, changed)
}

// DiffAgainstBaseline compares current secret data against a saved baseline.
func DiffAgainstBaseline(b *Baseline, current map[string]string) (*BaselineDiffResult, error) {
	if b == nil {
		return nil, fmt.Errorf("baseline is nil")
	}
	changes := DiffSecrets(b.Data, current)
	return &BaselineDiffResult{
		Environment: b.Environment,
		Path:        b.Path,
		Changes:     changes,
	}, nil
}
