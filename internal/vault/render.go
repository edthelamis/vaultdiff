package vault

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// RenderOptions controls output formatting.
type RenderOptions struct {
	Color   bool
	Verbose bool
}

// RenderDiff writes a human-readable diff to the provided writer.
func RenderDiff(w io.Writer, path string, diff DiffResult, opts RenderOptions) {
	fmt.Fprintf(w, "Secret path: %s\n", path)

	if !diff.HasChanges() {
		fmt.Fprintln(w, "  No changes detected.")
		return
	}

	for _, k := range SortedKeys(diff.Added) {
		line := fmt.Sprintf("  + %s = %s", k, diff.Added[k])
		fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))
	}

	for _, k := range SortedKeys(diff.Removed) {
		line := fmt.Sprintf("  - %s = %s", k, diff.Removed[k])
		fmt.Fprintln(w, colorize(line, colorRed, opts.Color))
	}

	for _, k := range sortedChangedKeys(diff.Changed) {
		pair := diff.Changed[k]
		oldLine := fmt.Sprintf("  ~ %s: %s -> %s", k, pair[0], pair[1])
		fmt.Fprintln(w, colorize(oldLine, colorYellow, opts.Color))
	}

	if opts.Verbose {
		for _, k := range SortedKeys(diff.Unchanged) {
			line := fmt.Sprintf("    %s = %s", k, diff.Unchanged[k])
			fmt.Fprintln(w, colorize(line, colorGray, opts.Color))
		}
	}
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}

func sortedChangedKeys(m map[string][2]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// reuse strings.Compare via sort
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if strings.Compare(keys[i], keys[j]) > 0 {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
