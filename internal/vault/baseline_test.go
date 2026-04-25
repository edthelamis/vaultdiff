package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBaseline_CopiesData(t *testing.T) {
	data := map[string]string{"key1": "val1", "key2": "val2"}
	b := NewBaseline("prod", "secret/app", data)
	if b.Environment != "prod" {
		t.Errorf("expected environment prod, got %s", b.Environment)
	}
	if b.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", b.Path)
	}
	if b.Data["key1"] != "val1" {
		t.Errorf("expected key1=val1")
	}
	// Mutating original should not affect baseline
	data["key1"] = "mutated"
	if b.Data["key1"] != "val1" {
		t.Error("baseline data was mutated via original map")
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	b := &Baseline{
		Environment: "staging",
		Path:        "secret/svc",
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		Data:        map[string]string{"db_pass": "secret"},
	}
	if err := SaveBaseline(b, path); err != nil {
		t.Fatalf("SaveBaseline failed: %v", err)
	}
	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline failed: %v", err)
	}
	if loaded.Environment != b.Environment {
		t.Errorf("environment mismatch: got %s", loaded.Environment)
	}
	if loaded.Data["db_pass"] != "secret" {
		t.Errorf("data mismatch")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveBaseline_NilBaseline(t *testing.T) {
	err := SaveBaseline(nil, "/tmp/noop.json")
	if err == nil {
		t.Error("expected error for nil baseline")
	}
}

func TestSaveBaseline_InvalidPath(t *testing.T) {
	b := NewBaseline("dev", "secret/x", map[string]string{})
	err := SaveBaseline(b, "/nonexistent/dir/baseline.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)
	_, err := LoadBaseline(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
