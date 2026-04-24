package vault

import "sort"

// DiffResult represents the diff between two secret versions.
type DiffResult struct {
	Added    map[string]string
	Removed  map[string]string
	Changed  map[string][2]string // key -> [oldValue, newValue]
	Unchanged map[string]string
}

// DiffSecrets compares two maps of secret key-value pairs and returns a DiffResult.
func DiffSecrets(old, new map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for k, newVal := range new {
		oldVal, exists := old[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = [2]string{oldVal, newVal}
		} else {
			result.Unchanged[k] = newVal
		}
	}

	for k, oldVal := range old {
		if _, exists := new[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// HasChanges returns true if there are any additions, removals, or changes.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// SortedKeys returns a sorted slice of keys from a map.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
