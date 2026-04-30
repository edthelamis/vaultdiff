package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestArchive_AddAndFilterByPath(t *testing.T) {
	a := NewArchive()
	a.Add("secret/app/prod", map[string]interface{}{"key": "val"}, "v1")
	a.Add("secret/app/dev", map[string]interface{}{"key": "devval"}, "")
	a.Add("secret/app/prod", map[string]interface{}{"key": "val2"}, "v2")

	results := a.FilterByPath("secret/app/prod")
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
	if results[0].Label != "v1" || results[1].Label != "v2" {
		t.Errorf("unexpected labels: %v, %v", results[0].Label, results[1].Label)
	}
}

func TestArchive_Add_CopiesData(t *testing.T) {
	a := NewArchive()
	data := map[string]interface{}{"token": "secret"}
	a.Add("secret/svc", data, "")
	data["token"] = "changed"

	if a.Entries[0].Data["token"] != "secret" {
		t.Error("archive entry should not reflect mutation of original map")
	}
}

func TestArchive_Summary(t *testing.T) {
	a := NewArchive()
	a.Add("secret/a", map[string]interface{}{"k": "v"}, "")
	a.Add("secret/a", map[string]interface{}{"k": "v2"}, "")
	a.Add("secret/b", map[string]interface{}{"x": "y"}, "")

	summary := a.Summary()
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	for _, want := range []string{"3 entries", "2 path", "secret/a", "secret/b"} {
		if !contains(summary, want) {
			t.Errorf("summary missing %q: %s", want, summary)
		}
	}
}

func TestSaveAndLoadArchive_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	a := NewArchive()
	a.Add("secret/x", map[string]interface{}{"foo": "bar"}, "label1")

	if err := SaveArchive(a, path); err != nil {
		t.Fatalf("SaveArchive: %v", err)
	}
	loaded, err := LoadArchive(path)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Label != "label1" {
		t.Errorf("expected label1, got %s", loaded.Entries[0].Label)
	}
}

func TestLoadArchive_MissingFile(t *testing.T) {
	_, err := LoadArchive("/nonexistent/archive.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveArchive_NilArchive(t *testing.T) {
	err := SaveArchive(nil, "/tmp/x.json")
	if err == nil {
		t.Fatal("expected error for nil archive")
	}
}

func TestArchive_ArchivedAtIsSet(t *testing.T) {
	before := time.Now().UTC()
	a := NewArchive()
	a.Add("secret/ts", map[string]interface{}{"k": "v"}, "")
	after := time.Now().UTC()

	at := a.Entries[0].ArchivedAt
	if at.Before(before) || at.After(after) {
		t.Errorf("ArchivedAt %v not in expected range [%v, %v]", at, before, after)
	}
}

func TestLoadArchiveStore_MissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	opts := ArchiveStoreOptions{FilePath: filepath.Join(dir, "missing.json")}
	a, err := LoadArchiveStore(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.Entries) != 0 {
		t.Errorf("expected empty archive, got %d entries", len(a.Entries))
	}
}

func TestSaveArchiveStore_EmptyPath(t *testing.T) {
	a := NewArchive()
	err := SaveArchiveStore(a, ArchiveStoreOptions{})
	if err == nil {
		t.Fatal("expected error for empty file path")
	}
}

// contains is a helper to check substring presence.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}

var _ = os.WriteFile // suppress unused import
