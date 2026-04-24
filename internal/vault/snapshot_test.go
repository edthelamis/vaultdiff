package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTakeSnapshot_CopiesData(t *testing.T) {
	original := map[string]string{"key1": "val1", "key2": "val2"}
	s := TakeSnapshot("prod", "secret/app", 3, original)

	original["key1"] = "mutated"
	if s.Data["key1"] != "val1" {
		t.Errorf("expected snapshot data to be independent of original map")
	}
	if s.Environment != "prod" || s.Path != "secret/app" || s.Version != 3 {
		t.Errorf("unexpected snapshot metadata: %+v", s)
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := TakeSnapshot("staging", "secret/db", 7, map[string]string{
		"host": "localhost",
		"port": "5432",
	})
	original.CapturedAt = original.CapturedAt.Truncate(time.Second)

	if err := SaveSnapshot(original, path); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if loaded.Environment != original.Environment {
		t.Errorf("environment mismatch: got %q want %q", loaded.Environment, original.Environment)
	}
	if loaded.Version != original.Version {
		t.Errorf("version mismatch: got %d want %d", loaded.Version, original.Version)
	}
	if loaded.Data["host"] != "localhost" || loaded.Data["port"] != "5432" {
		t.Errorf("data mismatch: %v", loaded.Data)
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	s := TakeSnapshot("dev", "secret/app", 1, map[string]string{})
	err := SaveSnapshot(s, "/nonexistent_dir/snap.json")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
