package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportFormat represents supported export formats.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportAuditLog writes the audit log to w in the specified format.
func ExportAuditLog(log *AuditLog, format ExportFormat, w io.Writer) error {
	switch format {
	case FormatJSON:
		return exportJSON(log, w)
	case FormatCSV:
		return exportCSV(log, w)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportJSON(log *AuditLog, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log.Entries)
}

func exportCSV(log *AuditLog, w io.Writer) error {
	_, err := fmt.Fprintln(w, "timestamp,environment,path,version_a,version_b,key,status,old_value,new_value")
	if err != nil {
		return err
	}
	for _, e := range log.Entries {
		ts := e.Timestamp.Format("2006-01-02T15:04:05Z")
		if len(e.Changes) == 0 {
			_, err = fmt.Fprintf(w, "%s,%s,%s,%d,%d,,,,\n",
				ts, e.Environment, e.Path, e.VersionA, e.VersionB)
			if err != nil {
				return err
			}
			continue
		}
		for _, c := range e.Changes {
			line := fmt.Sprintf("%s,%s,%s,%d,%d,%s,%s,%s,%s\n",
				ts,
				e.Environment,
				e.Path,
				e.VersionA,
				e.VersionB,
				c.Key,
				string(c.Status),
				escapeCSV(c.OldValue),
				escapeCSV(c.NewValue),
			)
			if _, err = fmt.Fprint(w, line); err != nil {
				return err
			}
		}
	}
	return nil
}

// escapeCSV replaces commas with semicolons and wraps values containing
// newlines or double-quotes in double-quotes, escaping internal quotes.
func escapeCSV(value string) string {
	value = strings.ReplaceAll(value, ",", ";")
	if strings.ContainsAny(value, "\"\n") {
		value = `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return value
}
