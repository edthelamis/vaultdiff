package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestVaultServer(t *testing.T, handler http.HandlerFunc) *vaultapi.Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	cfg := vaultapi.DefaultConfig()
	cfg.Address = ts.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("creating vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestGetSecretVersion_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{"API_KEY": "abc123", "DB_PASS": "secret"},
			"metadata": map[string]interface{}{"version": float64(3)},
		},
	}

	client := newTestVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	})

	sv, err := GetSecretVersion(context.Background(), client, "secret", "myapp/config", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sv.Version != 3 {
		t.Errorf("expected version 3, got %d", sv.Version)
	}
	if sv.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", sv.Data["API_KEY"])
	}
}

func TestGetSecretVersion_NotFound(t *testing.T) {
	client := newTestVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := GetSecretVersion(context.Background(), client, "secret", "missing/path", 0)
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
