package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api     *vaultapi.Client
	Address string
}

// NewClient creates a new Vault client using the provided address and token.
// Falls back to VAULT_ADDR and VAULT_TOKEN environment variables if empty.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()

	if address != "" {
		cfg.Address = address
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	} else {
		return nil, fmt.Errorf("vault address not provided and VAULT_ADDR is not set")
	}

	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	} else if t := os.Getenv("VAULT_TOKEN"); t != "" {
		client.SetToken(t)
	} else {
		return nil, fmt.Errorf("vault token not provided and VAULT_TOKEN is not set")
	}

	return &Client{
		api:     client,
		Address: cfg.Address,
	}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// Pass version 0 to read the latest version.
func (c *Client) ReadSecretVersion(mount, path string, version int) (map[string]interface{}, error) {
	var versionParam map[string][]string
	if version > 0 {
		versionParam = map[string][]string{
			"version": {fmt.Sprintf("%d", version)},
		}
	}

	secret, err := c.api.KVv2(mount).GetVersion(nil, path, version)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret %s/%s (version %d): %w", mount, path, version, err)
	}
	_ = versionParam

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %s/%s (version %d) not found or has no data", mount, path, version)
	}

	return secret.Data, nil
}
