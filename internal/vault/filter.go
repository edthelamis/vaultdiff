package vault

import "strings"

// FilterOptions controls which diff entries are included in output.
type FilterOptions struct {
	// ChangeTypes limits results to specific change types (added, removed, changed, unchanged).
	// If empty, all types are included.
	ChangeTypes []string

	// KeyPrefix filters entries to only those whose key starts with the given prefix.
	KeyPrefix string

	// ExcludeKeys is a list of exact key names to omit from results.
	ExcludeKeys []string
}

// FilterDiff returns a filtered subset of DiffEntry values based on the provided options.
func FilterDiff(entries []DiffEntry, opts FilterOptions) []DiffEntry {
	allowedTypes := make(map[string]bool)
	for _, t := range opts.ChangeTypes {
		allowedTypes[strings.ToLower(t)] = true
	}

	excluded := make(map[string]bool)
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	var result []DiffEntry
	for _, e := range entries {
		if len(allowedTypes) > 0 && !allowedTypes[strings.ToLower(string(e.Type))] {
			continue
		}
		if opts.KeyPrefix != "" && !strings.HasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		if excluded[e.Key] {
			continue
		}
		result = append(result, e)
	}
	return result
}
