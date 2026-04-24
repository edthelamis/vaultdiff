package vault

import (
	"fmt"
	"io"
	"sort"
)

// ChangeType constants used for colorization and display.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// RenderDiff writes a human-readable diff of the provided entries to w.
// If color is true, ANSI color codes are used.
func RenderDiff(w io.Writer, entries []DiffEntry, color bool) {
	for _, e := range entries {
		line := formatEntry(e, color)
		fmt.Fprintln(w, line)
	}
}

func formatEntry(e DiffEntry, color bool) string {
	switch e.Type {
	case ChangeTypeAdded:
		return colorize(fmt.Sprintf("+ %s = %q", e.Key, e.NewValue), colorGreen, color)
	case ChangeTypeRemoved:
		return colorize(fmt.Sprintf("- %s = %q", e.Key, e.OldValue), colorRed, color)
	case ChangeTypeChanged:
		return colorize(fmt.Sprintf("~ %s: %q -> %q", e.Key, e.OldValue, e.NewValue), colorYellow, color)
	case ChangeTypeUnchanged:
		return colorize(fmt.Sprintf("  %s = %q", e.Key, e.NewValue), colorGray, color)
	default:
		return fmt.Sprintf("? %s", e.Key)
	}
}

func colorize(s, code string, enabled bool) string {
	if !enabled {
		return s
	}
	return code + s + colorReset
}

func sortedChangedKeys(m map[string]DiffEntry) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
