package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleHistoryChanges() []DiffEntry {
	return []DiffEntry{
		{Key: "db_pass", ChangeType: Changed, OldValue: "old", NewValue: "new"},
		{Key: "api_key", ChangeType: Added, OldValue: "", NewValue: "abc"},
	}
}

func TestHistory_Record(t *testing.T) {
	h := NewHistory()
	h.Record("prod", "secret/app", 1, 2, sampleHistoryChanges())
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	e := h.Entries[0]
	if e.Environment != "prod" || e.Path != "secret/app" {
		t.Errorf("unexpected entry values: %+v", e)
	}
	if e.FromVersion != 1 || e.ToVersion != 2 {
		t.Errorf("unexpected versions: from=%d to=%d", e.FromVersion, e.ToVersion)
	}
}

func TestHistory_Filter(t *testing.T) {
	h := NewHistory()
	h.Record("prod", "secret/app", 1, 2, sampleHistoryChanges())
	h.Record("staging", "secret/app", 1, 2, sampleHistoryChanges())
	h.Record("prod", "secret/db", 3, 4, sampleHistoryChanges())

	prodEntries := h.Filter("prod", "")
	if len(prodEntries) != 2 {
		t.Errorf("expected 2 prod entries, got %d", len(prodEntries))
	}

	specific := h.Filter("prod", "secret/db")
	if len(specific) != 1 {
		t.Errorf("expected 1 specific entry, got %d", len(specific))
	}

	all := h.Filter("", "")
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}

func TestHistory_Summary_Empty(t *testing.T) {
	h := NewHistory()
	if h.Summary() != "No history recorded." {
		t.Errorf("unexpected summary for empty history")
	}
}

func TestHistory_Summary_WithEntries(t *testing.T) {
	h := NewHistory()
	h.Record("prod", "secret/app", 1, 2, sampleHistoryChanges())
	summary := h.Summary()
	if !strings.Contains(summary, "prod") || !strings.Contains(summary, "secret/app") {
		t.Errorf("summary missing expected fields: %s", summary)
	}
	if !strings.Contains(summary, "2 changes") {
		t.Errorf("summary should report 2 changes, got: %s", summary)
	}
}

func TestSaveAndLoadHistory_RoundTrip(t *testing.T) {
	h := NewHistory()
	h.Record("dev", "secret/svc", 2, 3, sampleHistoryChanges())

	tmp := filepath.Join(t.TempDir(), "history.json")
	if err := SaveHistory(h, tmp); err != nil {
		t.Fatalf("SaveHistory error: %v", err)
	}
	loaded, err := LoadHistory(tmp)
	if err != nil {
		t.Fatalf("LoadHistory error: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Environment != "dev" {
		t.Errorf("unexpected environment: %s", loaded.Entries[0].Environment)
	}
}

func TestLoadHistory_MissingFile(t *testing.T) {
	_, err := LoadHistory("/nonexistent/path/history.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveHistory_NilHistory(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "h.json")
	if err := SaveHistory(nil, tmp); err == nil {
		t.Fatal("expected error for nil history")
	}
}

func TestSaveHistory_InvalidPath(t *testing.T) {
	h := NewHistory()
	// Use a file as a directory component to force failure
	f, _ := os.CreateTemp("", "file")
	f.Close()
	err := SaveHistory(h, filepath.Join(f.Name(), "sub", "h.json"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
