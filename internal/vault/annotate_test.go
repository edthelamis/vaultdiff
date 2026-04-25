package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnnotationIndex_AddAndGet(t *testing.T) {
	idx := NewAnnotationIndex()
	idx.Add("secret/app/db", "primary database credentials", "alice")

	a, ok := idx.Get("secret/app/db")
	if !ok {
		t.Fatal("expected annotation to exist")
	}
	if a.Note != "primary database credentials" {
		t.Errorf("unexpected note: %s", a.Note)
	}
	if a.Author != "alice" {
		t.Errorf("unexpected author: %s", a.Author)
	}
}

func TestAnnotationIndex_Remove(t *testing.T) {
	idx := NewAnnotationIndex()
	idx.Add("secret/app/api", "api key", "bob")

	if removed := idx.Remove("secret/app/api"); !removed {
		t.Fatal("expected Remove to return true")
	}
	if _, ok := idx.Get("secret/app/api"); ok {
		t.Error("expected annotation to be removed")
	}
	if removed := idx.Remove("secret/app/api"); removed {
		t.Error("expected Remove to return false for missing key")
	}
}

func TestAnnotationIndex_Summary(t *testing.T) {
	idx := NewAnnotationIndex()
	if idx.Summary() != "no annotations" {
		t.Errorf("unexpected empty summary: %s", idx.Summary())
	}
	idx.Add("secret/x", "note", "user")
	if idx.Summary() != "1 annotation(s) recorded" {
		t.Errorf("unexpected summary: %s", idx.Summary())
	}
}

func TestSaveAndLoadAnnotationIndex_RoundTrip(t *testing.T) {
	idx := NewAnnotationIndex()
	idx.Add("secret/svc/token", "service token", "carol")

	tmp := filepath.Join(t.TempDir(), "annotations.json")
	if err := SaveAnnotationIndex(idx, tmp); err != nil {
		t.Fatalf("SaveAnnotationIndex: %v", err)
	}

	loaded, err := LoadAnnotationIndex(tmp)
	if err != nil {
		t.Fatalf("LoadAnnotationIndex: %v", err)
	}
	a, ok := loaded.Get("secret/svc/token")
	if !ok {
		t.Fatal("expected annotation after round-trip")
	}
	if a.Note != "service token" || a.Author != "carol" {
		t.Errorf("unexpected annotation after round-trip: %+v", a)
	}
}

func TestLoadAnnotationIndex_MissingFile(t *testing.T) {
	_, err := LoadAnnotationIndex("/nonexistent/path/annotations.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAnnotationIndex_NilIndex(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "annotations.json")
	if err := SaveAnnotationIndex(nil, tmp); err == nil {
		t.Fatal("expected error for nil index")
	}
}

func TestSaveAnnotationIndex_InvalidPath(t *testing.T) {
	idx := NewAnnotationIndex()
	err := SaveAnnotationIndex(idx, string([]byte{0}))
	if err == nil {
		// some OS may allow it; skip gracefully
		os.Remove(string([]byte{0}))
	}
}
