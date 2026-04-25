package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRollbackVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// Serve version 1 data.
	mux.HandleFunc("/v1/secret/data/myapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"data": map[string]string{"key": "old-value"},
					"metadata": map[string]interface{}{"version": 1},
				},
			})
			return
		}
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{}})
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return httptest.NewServer(mux)
}

func TestRollbackSecret_DryRun(t *testing.T) {
	srv := newRollbackVaultServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	opts := DefaultRollbackOptions()
	opts.DryRun = true

	result, err := RollbackSecret(context.Background(), client, "myapp", 1, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RolledBack {
		t.Error("expected RolledBack=false in dry-run mode")
	}
	if !result.DryRun {
		t.Error("expected DryRun=true on result")
	}
	if result.ToVersion != 1 {
		t.Errorf("expected ToVersion=1, got %d", result.ToVersion)
	}
}

func TestRollbackSecret_NilClient(t *testing.T) {
	_, err := RollbackSecret(context.Background(), nil, "myapp", 1, DefaultRollbackOptions())
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRollbackSecret_EmptyPath(t *testing.T) {
	srv := newRollbackVaultServer(t)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token")
	_, err := RollbackSecret(context.Background(), client, "", 1, DefaultRollbackOptions())
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRollbackSecret_InvalidVersion(t *testing.T) {
	srv := newRollbackVaultServer(t)
	defer srv.Close()

	client, _ := NewClient(srv.URL, "test-token")
	_, err := RollbackSecret(context.Background(), client, "myapp", 0, DefaultRollbackOptions())
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestDefaultRollbackOptions(t *testing.T) {
	opts := DefaultRollbackOptions()
	if opts.MountPath != "secret" {
		t.Errorf("expected default MountPath=secret, got %q", opts.MountPath)
	}
	if opts.DryRun {
		t.Error("expected DryRun=false by default")
	}
}
