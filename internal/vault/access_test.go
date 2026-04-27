package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAccessLog_RecordAndFilter(t *testing.T) {
	log := NewAccessLog()
	log.Record("secret/app", "read", "alice")
	log.Record("secret/db", "write", "bob")
	log.Record("secret/app", "delete", "alice")

	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log.Entries))
	}

	byAlice := log.FilterByActor("alice")
	if len(byAlice) != 2 {
		t.Errorf("expected 2 entries for alice, got %d", len(byAlice))
	}

	byPath := log.FilterByPath("secret/app")
	if len(byPath) != 2 {
		t.Errorf("expected 2 entries for secret/app, got %d", len(byPath))
	}
}

func TestAccessLog_Summary_Empty(t *testing.T) {
	log := NewAccessLog()
	got := log.Summary()
	if got != "no access events recorded" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestAccessLog_Summary_WithEntries(t *testing.T) {
	log := NewAccessLog()
	log.Record("secret/app", "read", "alice")
	log.Record("secret/app", "write", "bob")
	log.Record("secret/db", "read", "alice")

	summary := log.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	for _, want := range []string{"3 total", "secret/app", "secret/db"} {
		if !containsStr(summary, want) {
			t.Errorf("summary missing %q: %s", want, summary)
		}
	}
}

func TestSaveAndLoadAccessLog_RoundTrip(t *testing.T) {
	log := NewAccessLog()
	log.Record("secret/app", "read", "alice")
	log.Record("secret/db", "delete", "bob")

	dir := t.TempDir()
	path := filepath.Join(dir, "access.json")

	if err := SaveAccessLog(log, path); err != nil {
		t.Fatalf("SaveAccessLog: %v", err)
	}

	loaded, err := LoadAccessLog(path)
	if err != nil {
		t.Fatalf("LoadAccessLog: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Actor != "alice" {
		t.Errorf("expected actor alice, got %s", loaded.Entries[0].Actor)
	}
}

func TestLoadAccessLog_MissingFile(t *testing.T) {
	_, err := LoadAccessLog("/nonexistent/access.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveAccessLog_InvalidPath(t *testing.T) {
	log := NewAccessLog()
	log.Record("secret/app", "read", "alice")
	err := SaveAccessLog(log, "/nonexistent/dir/access.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestSaveAccessLog_Nil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "access.json")
	err := SaveAccessLog(nil, path)
	if err == nil {
		t.Error("expected error for nil log")
	}
}

// containsStr checks if s contains substr.
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

var _ = os.DevNull // suppress unused import
