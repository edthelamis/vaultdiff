package vault

import (
	"os"
	"path/filepath"
	"testing"
)

var testSecret = []byte("super-secret-hmac-key")

var sampleSecretData = map[string]interface{}{
	"db_password": "hunter2",
	"api_key":     "abc123",
	"debug":       false,
}

func TestSignSecrets_ProducesSignatures(t *testing.T) {
	rec, err := SignSecrets("secret/myapp", sampleSecretData, SignOptions{Secret: testSecret})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", rec.Path)
	}
	if len(rec.Signatures) != len(sampleSecretData) {
		t.Errorf("expected %d signatures, got %d", len(sampleSecretData), len(rec.Signatures))
	}
}

func TestSignSecrets_EmptySecret(t *testing.T) {
	_, err := SignSecrets("secret/myapp", sampleSecretData, SignOptions{})
	if err == nil {
		t.Fatal("expected error for empty HMAC secret")
	}
}

func TestSignSecrets_EmptyPath(t *testing.T) {
	_, err := SignSecrets("", sampleSecretData, SignOptions{Secret: testSecret})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestVerifySecrets_Valid(t *testing.T) {
	rec, _ := SignSecrets("secret/myapp", sampleSecretData, SignOptions{Secret: testSecret})
	failed, err := VerifySecrets(sampleSecretData, rec, SignOptions{Secret: testSecret})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(failed) != 0 {
		t.Errorf("expected no failures, got: %v", failed)
	}
}

func TestVerifySecrets_TamperedValue(t *testing.T) {
	rec, _ := SignSecrets("secret/myapp", sampleSecretData, SignOptions{Secret: testSecret})
	tampered := map[string]interface{}{
		"db_password": "TAMPERED",
		"api_key":     "abc123",
		"debug":       false,
	}
	failed, err := VerifySecrets(tampered, rec, SignOptions{Secret: testSecret})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(failed) != 1 || failed[0] != "db_password" {
		t.Errorf("expected db_password to fail, got: %v", failed)
	}
}

func TestVerifySecrets_MissingKey(t *testing.T) {
	rec, _ := SignSecrets("secret/myapp", sampleSecretData, SignOptions{Secret: testSecret})
	partial := map[string]interface{}{
		"api_key": "abc123",
	}
	failed, err := VerifySecrets(partial, rec, SignOptions{Secret: testSecret})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(failed) == 0 {
		t.Error("expected failures for missing keys")
	}
}

func TestVerifySecrets_NilRecord(t *testing.T) {
	_, err := VerifySecrets(sampleSecretData, nil, SignOptions{Secret: testSecret})
	if err == nil {
		t.Fatal("expected error for nil record")
	}
}

func TestSaveAndLoadSignatureRecord_RoundTrip(t *testing.T) {
	rec, _ := SignSecrets("secret/myapp", sampleSecretData, SignOptions{Secret: testSecret})
	dir := t.TempDir()
	path := filepath.Join(dir, "sig.json")

	if err := SaveSignatureRecord(path, rec); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadSignatureRecord(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Path != rec.Path {
		t.Errorf("path mismatch: got %s", loaded.Path)
	}
	if len(loaded.Signatures) != len(rec.Signatures) {
		t.Errorf("signature count mismatch")
	}
}

func TestLoadSignatureRecord_MissingFile(t *testing.T) {
	_, err := LoadSignatureRecord("/nonexistent/path/sig.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveSignatureRecord_NilRecord(t *testing.T) {
	dir := t.TempDir()
	err := SaveSignatureRecord(filepath.Join(dir, "sig.json"), nil)
	if err == nil {
		t.Fatal("expected error for nil record")
	}
}

func TestSaveSignatureRecord_InvalidPath(t *testing.T) {
	rec, _ := SignSecrets("secret/myapp", sampleSecretData, SignOptions{Secret: testSecret})
	err := SaveSignatureRecord("/dev/null/invalid/sig.json", rec)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
	_ = os.Remove("/dev/null/invalid/sig.json")
}
