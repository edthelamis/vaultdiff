package vault

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// SignatureRecord holds a signed snapshot of secret keys and their HMAC signatures.
type SignatureRecord struct {
	Path      string            `json:"path"`
	Timestamp time.Time         `json:"timestamp"`
	Signatures map[string]string `json:"signatures"`
}

// SignOptions configures signing behaviour.
type SignOptions struct {
	Secret []byte // HMAC secret key
}

// SignSecrets creates a SignatureRecord for the given secret data using HMAC-SHA256.
func SignSecrets(path string, data map[string]interface{}, opts SignOptions) (*SignatureRecord, error) {
	if len(opts.Secret) == 0 {
		return nil, fmt.Errorf("sign: HMAC secret must not be empty")
	}
	if path == "" {
		return nil, fmt.Errorf("sign: path must not be empty")
	}

	sigs := make(map[string]string, len(data))
	for k, v := range data {
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("sign: marshal key %q: %w", k, err)
		}
		mac := hmac.New(sha256.New, opts.Secret)
		mac.Write(raw)
		sigs[k] = hex.EncodeToString(mac.Sum(nil))
	}

	return &SignatureRecord{
		Path:       path,
		Timestamp:  time.Now().UTC(),
		Signatures: sigs,
	}, nil
}

// VerifySecrets checks the given data against a SignatureRecord.
// It returns a list of keys that failed verification and any error.
func VerifySecrets(data map[string]interface{}, record *SignatureRecord, opts SignOptions) ([]string, error) {
	if record == nil {
		return nil, fmt.Errorf("verify: signature record is nil")
	}
	if len(opts.Secret) == 0 {
		return nil, fmt.Errorf("verify: HMAC secret must not be empty")
	}

	var failed []string

	// Check every key in the record.
	keys := make([]string, 0, len(record.Signatures))
	for k := range record.Signatures {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		expected := record.Signatures[k]
		v, ok := data[k]
		if !ok {
			failed = append(failed, k)
			continue
		}
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("verify: marshal key %q: %w", k, err)
		}
		mac := hmac.New(sha256.New, opts.Secret)
		mac.Write(raw)
		actual := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(actual), []byte(expected)) {
			failed = append(failed, k)
		}
	}

	return failed, nil
}
