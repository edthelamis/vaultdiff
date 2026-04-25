package vault

import "fmt"

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy string

const (
	MergeStrategySource      MergeStrategy = "source"       // source always wins
	MergeStrategyDestination MergeStrategy = "destination"  // destination always wins
	MergeStrategyError       MergeStrategy = "error"        // conflict returns an error
)

// MergeOptions configures the behaviour of MergeSecrets.
type MergeOptions struct {
	Strategy    MergeStrategy
	DryRun      bool
	ExcludeKeys []string
}

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Merged    map[string]string
	Added     []string
	Overwritten []string
	Skipped   []string
	Conflicts []string
}

// Summary returns a human-readable summary of the merge result.
func (r *MergeResult) Summary() string {
	return fmt.Sprintf(
		"merged=%d added=%d overwritten=%d skipped=%d conflicts=%d",
		len(r.Merged), len(r.Added), len(r.Overwritten), len(r.Skipped), len(r.Conflicts),
	)
}

// MergeSecrets merges src into dst according to the provided options.
// It returns a MergeResult describing what changed, or an error when
// MergeStrategyError is used and conflicts are detected.
func MergeSecrets(dst, src map[string]string, opts MergeOptions) (*MergeResult, error) {
	if opts.Strategy == "" {
		opts.Strategy = MergeStrategySource
	}

	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	result := &MergeResult{
		Merged: make(map[string]string),
	}

	// Start with a copy of dst.
	for k, v := range dst {
		result.Merged[k] = v
	}

	for k, srcVal := range src {
		if excluded[k] {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		dstVal, exists := dst[k]
		if !exists {
			result.Added = append(result.Added, k)
			if !opts.DryRun {
				result.Merged[k] = srcVal
			}
			continue
		}

		if dstVal == srcVal {
			// No conflict, values are identical — nothing to do.
			continue
		}

		// Conflict: values differ.
		switch opts.Strategy {
		case MergeStrategySource:
			result.Overwritten = append(result.Overwritten, k)
			if !opts.DryRun {
				result.Merged[k] = srcVal
			}
		case MergeStrategyDestination:
			result.Skipped = append(result.Skipped, k)
		case MergeStrategyError:
			result.Conflicts = append(result.Conflicts, k)
		}
	}

	if len(result.Conflicts) > 0 {
		return result, fmt.Errorf("merge conflict on keys: %v", result.Conflicts)
	}
	return result, nil
}
