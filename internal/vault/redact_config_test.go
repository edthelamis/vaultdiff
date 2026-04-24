package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeRedactConfig(t *testing.T, cfg RedactConfig) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	p := filepath.Join(dir, "redact.json")
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadRedactConfig_Valid(t *testing.T) {
	path := writeRedactConfig(t, RedactConfig{
		KeyPatterns: []string{"*password*", "*token*"},
		Replacement: "***",
	})
	opts, err := LoadRedactConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.KeyPatterns) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(opts.KeyPatterns))
	}
	if opts.Replacement != "***" {
		t.Errorf("expected replacement '***', got %q", opts.Replacement)
	}
}

func TestLoadRedactConfig_MissingFile(t *testing.T) {
	opts, err := LoadRedactConfig("/nonexistent/path/redact.json")
	if err != nil {
		t.Fatalf("expected fallback to defaults, got error: %v", err)
	}
	if len(opts.KeyPatterns) == 0 {
		t.Error("expected default key patterns")
	}
}

func TestLoadRedactConfig_EmptyPath(t *testing.T) {
	opts, err := LoadRedactConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.KeyPatterns) == 0 {
		t.Error("expected default key patterns for empty path")
	}
}

func TestLoadRedactConfig_MissingPatterns(t *testing.T) {
	path := writeRedactConfig(t, RedactConfig{
		KeyPatterns: []string{},
		Replacement: "[REDACTED]",
	})
	_, err := LoadRedactConfig(path)
	if err == nil {
		t.Error("expected error for missing key_patterns")
	}
}

func TestLoadRedactConfig_DefaultReplacement(t *testing.T) {
	path := writeRedactConfig(t, RedactConfig{
		KeyPatterns: []string{"*secret*"},
		Replacement: "",
	})
	opts, err := LoadRedactConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Replacement != "[REDACTED]" {
		t.Errorf("expected default replacement, got %q", opts.Replacement)
	}
}
