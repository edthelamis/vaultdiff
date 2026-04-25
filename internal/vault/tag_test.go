package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTagIndex_AddAndFilterByTag(t *testing.T) {
	idx := make(TagIndex)
	idx.AddTag("secret/app/db", "env", "prod")
	idx.AddTag("secret/app/api", "env", "prod")
	idx.AddTag("secret/app/cache", "env", "staging")

	results := idx.FilterByTag("env", "prod")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "secret/app/api" || results[1] != "secret/app/db" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestTagIndex_RemoveTag(t *testing.T) {
	idx := make(TagIndex)
	idx.AddTag("secret/app/db", "env", "prod")
	idx.RemoveTag("secret/app/db", "env")

	if _, ok := idx["secret/app/db"]; ok {
		t.Error("expected path to be removed after last tag deleted")
	}
}

func TestTagIndex_RemoveTag_Partial(t *testing.T) {
	idx := make(TagIndex)
	idx.AddTag("secret/app/db", "env", "prod")
	idx.AddTag("secret/app/db", "team", "platform")
	idx.RemoveTag("secret/app/db", "env")

	tags, ok := idx["secret/app/db"]
	if !ok {
		t.Fatal("expected path to remain after partial tag removal")
	}
	if _, hasEnv := tags["env"]; hasEnv {
		t.Error("expected 'env' tag to be removed")
	}
	if tags["team"] != "platform" {
		t.Error("expected 'team' tag to remain")
	}
}

func TestTagIndex_Summary_Empty(t *testing.T) {
	idx := make(TagIndex)
	got := idx.Summary()
	if got != "no tagged secrets" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestTagIndex_Summary_WithEntries(t *testing.T) {
	idx := make(TagIndex)
	idx.AddTag("secret/app/db", "env", "prod")
	summary := idx.Summary()
	if summary == "" || summary == "no tagged secrets" {
		t.Errorf("expected non-empty summary, got: %q", summary)
	}
}

func TestSaveAndLoadTagIndex_RoundTrip(t *testing.T) {
	idx := make(TagIndex)
	idx.AddTag("secret/app/db", "env", "prod")
	idx.AddTag("secret/app/api", "team", "platform")

	tmp := filepath.Join(t.TempDir(), "tags.json")
	if err := SaveTagIndex(idx, tmp); err != nil {
		t.Fatalf("SaveTagIndex: %v", err)
	}

	loaded, err := LoadTagIndex(tmp)
	if err != nil {
		t.Fatalf("LoadTagIndex: %v", err)
	}
	if loaded["secret/app/db"]["env"] != "prod" {
		t.Errorf("expected prod, got %q", loaded["secret/app/db"]["env"])
	}
	if loaded["secret/app/api"]["team"] != "platform" {
		t.Errorf("expected platform, got %q", loaded["secret/app/api"]["team"])
	}
}

func TestLoadTagIndex_MissingFile(t *testing.T) {
	_, err := LoadTagIndex("/nonexistent/tags.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveTagIndex_NilIndex(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "tags.json")
	err := SaveTagIndex(nil, tmp)
	if err == nil {
		t.Error("expected error for nil index")
	}
}

func TestSaveTagIndex_InvalidPath(t *testing.T) {
	idx := make(TagIndex)
	idx.AddTag("secret/app/db", "env", "prod")
	err := SaveTagIndex(idx, string([]byte{0}))
	if err == nil {
		// some systems may not error on null-byte paths, skip
		t.Skip("platform did not reject invalid path")
	}
	_ = os.Remove(string([]byte{0}))
}
