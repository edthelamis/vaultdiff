package vault

import (
	"fmt"
	"sort"
	"strings"
)

// TagSet holds a map of key-value tags associated with a secret path.
type TagSet map[string]string

// TaggedSecret associates a Vault secret path with a set of tags.
type TaggedSecret struct {
	Path string  `json:"path"`
	Tags TagSet  `json:"tags"`
}

// TagIndex maps secret paths to their tag sets.
type TagIndex map[string]TagSet

// AddTag adds or updates a tag on the given path in the index.
func (idx TagIndex) AddTag(path, key, value string) {
	if _, ok := idx[path]; !ok {
		idx[path] = make(TagSet)
	}
	idx[path][key] = value
}

// RemoveTag removes a tag key from the given path.
func (idx TagIndex) RemoveTag(path, key string) {
	if tags, ok := idx[path]; ok {
		delete(tags, key)
		if len(tags) == 0 {
			delete(idx, path)
		}
	}
}

// FilterByTag returns all paths whose tags contain the given key=value pair.
func (idx TagIndex) FilterByTag(key, value string) []string {
	var results []string
	for path, tags := range idx {
		if v, ok := tags[key]; ok && v == value {
			results = append(results, path)
		}
	}
	sort.Strings(results)
	return results
}

// Summary returns a human-readable summary of all tagged paths.
func (idx TagIndex) Summary() string {
	if len(idx) == 0 {
		return "no tagged secrets"
	}
	var sb strings.Builder
	paths := make([]string, 0, len(idx))
	for p := range idx {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		tags := idx[p]
		keys := make([]string, 0, len(tags))
		for k := range tags {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s=%s", k, tags[k]))
		}
		fmt.Fprintf(&sb, "%s [%s]\n", p, strings.Join(parts, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
