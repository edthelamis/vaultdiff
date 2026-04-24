package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func writePolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.json")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writing policy file: %v", err)
	}
	return p
}

func TestLoadPolicy_Valid(t *testing.T) {
	path := writePolicyFile(t, `{
		"name": "test-policy",
		"rules": [
			{"key_pattern": "db/", "allow_read": true, "allow_write": false, "allow_delete": false}
		]
	}`)
	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "test-policy" {
		t.Errorf("expected name test-policy, got %q", p.Name)
	}
	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
	if p.Rules[0].KeyPattern != "db/" {
		t.Errorf("unexpected key pattern: %q", p.Rules[0].KeyPattern)
	}
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	_, err := LoadPolicy("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadPolicy_InvalidJSON(t *testing.T) {
	path := writePolicyFile(t, `{not valid json}`)
	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadPolicy_MissingName(t *testing.T) {
	path := writePolicyFile(t, `{"rules": []}`)
	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoadPolicy_EmptyKeyPattern(t *testing.T) {
	path := writePolicyFile(t, `{
		"name": "bad",
		"rules": [{"key_pattern": "", "allow_read": true}]
	}`)
	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatal("expected error for empty key_pattern")
	}
}
