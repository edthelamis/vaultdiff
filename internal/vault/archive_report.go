package vault

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// ArchiveReport formats archive entries for display.
type ArchiveReport struct {
	Archive   *Archive
	FilterPath string // optional: only show entries for this path
}

// NewArchiveReport creates a report for the given archive.
func NewArchiveReport(a *Archive, filterPath string) *ArchiveReport {
	return &ArchiveReport{Archive: a, FilterPath: filterPath}
}

// Write renders the archive report to the provided writer.
func (r *ArchiveReport) Write(w io.Writer) error {
	if r.Archive == nil {
		return fmt.Errorf("archive is nil")
	}

	entries := r.Archive.Entries
	if r.FilterPath != "" {
		entries = r.Archive.FilterByPath(r.FilterPath)
	}

	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No archive entries found.")
		return err
	}

	// Sort by path then by archived time.
	sorted := make([]ArchiveEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Path != sorted[j].Path {
			return sorted[i].Path < sorted[j].Path
		}
		return sorted[i].ArchivedAt.Before(sorted[j].ArchivedAt)
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tLABEL\tARCHIVED AT\tKEYS")
	fmt.Fprintln(tw, "----\t-----\t-----------\t----")
	for _, e := range sorted {
		label := e.Label
		if label == "" {
			label = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n",
			e.Path,
			label,
			e.ArchivedAt.Format(time.RFC3339),
			len(e.Data),
		)
	}
	return tw.Flush()
}
