package vault

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildTestLog() *AuditLog {
	log := NewAuditLog()
	log.Record("production", "secret/myapp", 3, 4, []DiffEntry{
		{Key: "TOKEN", Status: StatusChanged, OldValue: "old_tok", NewValue: "new_tok"},
		{Key: "REMOVED_KEY", Status: StatusRemoved, OldValue: "gone", NewValue: ""},
	})
	return log
}

func TestExportAuditLog_JSON(t *testing.T) {
	log := buildTestLog()
	var buf bytes.Buffer
	if err := ExportAuditLog(log, FormatJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []AuditEntry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Environment != "production" {
		t.Errorf("unexpected environment: %q", entries[0].Environment)
	}
	if len(entries[0].Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(entries[0].Changes))
	}
}

func TestExportAuditLog_CSV(t *testing.T) {
	log := buildTestLog()
	var buf bytes.Buffer
	if err := ExportAuditLog(log, FormatCSV, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "timestamp,environment,path") {
		t.Error("CSV header missing")
	}
	if !strings.Contains(output, "TOKEN") {
		t.Error("CSV missing TOKEN key")
	}
	if !strings.Contains(output, "REMOVED_KEY") {
		t.Error("CSV missing REMOVED_KEY")
	}
}

func TestExportAuditLog_UnsupportedFormat(t *testing.T) {
	log := NewAuditLog()
	var buf bytes.Buffer
	err := ExportAuditLog(log, ExportFormat("xml"), &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExportAuditLog_CSV_NoChanges(t *testing.T) {
	log := NewAuditLog()
	log.Record("dev", "secret/empty", 1, 1, []DiffEntry{})
	var buf bytes.Buffer
	if err := ExportAuditLog(log, FormatCSV, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least header + 1 data row, got %d lines", len(lines))
	}
}
