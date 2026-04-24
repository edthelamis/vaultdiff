package vault_test

import (
	"os"
	"testing"

	"github.com/yourorg/vaultdiff/internal/vault"
)

func TestNewClient_MissingAddress(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")

	_, err := vault.NewClient("", "test-token")
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")

	_, err := vault.NewClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_FromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "root")

	client, err := vault.NewClient("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected address http://127.0.0.1:8200, got %s", client.Address)
	}
}

func TestNewClient_ExplicitValues(t *testing.T) {
	client, err := vault.NewClient("http://vault.example.com:8200", "my-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.Address != "http://vault.example.com:8200" {
		t.Errorf("expected address http://vault.example.com:8200, got %s", client.Address)
	}
}
