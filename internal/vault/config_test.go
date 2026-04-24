package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/vaultdiff/internal/vault"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultdiff-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	content := `
environments:
  - name: staging
    address: http://vault-staging:8200
    token: s.staging
    mount: secret
  - name: production
    address: http://vault-prod:8200
    token: s.prod
    mount: secret
`
	path := writeTemp(t, content)
	cfg, err := vault.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(cfg.Environments))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := vault.LoadConfig(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_DuplicateEnv(t *testing.T) {
	content := `
environments:
  - name: staging
    mount: secret
  - name: staging
    mount: secret
`
	path := writeTemp(t, content)
	_, err := vault.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for duplicate environment names, got nil")
	}
}

func TestGetEnvironment_NotFound(t *testing.T) {
	content := `
environments:
  - name: staging
    mount: secret
`
	path := writeTemp(t, content)
	cfg, err := vault.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = cfg.GetEnvironment("production")
	if err == nil {
		t.Fatal("expected error for unknown environment, got nil")
	}
}
