package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeLintConfig(t *testing.T, cfg LintConfig) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal lint config: %v", err)
	}
	dir := t.TempDir()
	p := filepath.Join(dir, "lint.json")
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatalf("failed to write lint config: %v", err)
	}
	return p
}

func TestLoadLintConfig_Valid(t *testing.T) {
	cfg := LintConfig{
		Name:  "prod-lint",
		Rules: []LintRule{{Name: "no-debug", Pattern: "debug_*", Target: "key", Message: "no debug keys"}},
	}
	p := writeLintConfig(t, cfg)
	loaded, err := LoadLintConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Name != "prod-lint" {
		t.Errorf("expected name prod-lint, got %s", loaded.Name)
	}
	if len(loaded.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(loaded.Rules))
	}
}

func TestLoadLintConfig_MissingFile(t *testing.T) {
	_, err := LoadLintConfig("/nonexistent/lint.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadLintConfig_EmptyPath(t *testing.T) {
	_, err := LoadLintConfig("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestLoadLintConfig_MissingName(t *testing.T) {
	cfg := LintConfig{
		Rules: []LintRule{{Name: "r1", Pattern: "*", Target: "key", Message: "msg"}},
	}
	p := writeLintConfig(t, cfg)
	_, err := LoadLintConfig(p)
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestLoadLintConfig_InvalidTarget(t *testing.T) {
	cfg := LintConfig{
		Name:  "bad-lint",
		Rules: []LintRule{{Name: "r1", Pattern: "*", Target: "invalid", Message: "msg"}},
	}
	p := writeLintConfig(t, cfg)
	_, err := LoadLintConfig(p)
	if err == nil {
		t.Error("expected error for invalid target")
	}
}

func TestLoadLintConfig_NoRules(t *testing.T) {
	cfg := LintConfig{Name: "empty-lint"}
	p := writeLintConfig(t, cfg)
	_, err := LoadLintConfig(p)
	if err == nil {
		t.Error("expected error for empty rules")
	}
}
