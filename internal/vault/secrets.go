package vault

import (
	"context"
	"fmt"
	"strconv"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretVersion holds the data and metadata for a specific version of a secret.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Metadata map[string]interface{}
}

// GetSecretVersion reads a specific version of a KV v2 secret from Vault.
// If version is 0, the latest version is returned.
func GetSecretVersion(ctx context.Context, client *vaultapi.Client, mount, path string, version int) (*SecretVersion, error) {
	kvPath := fmt.Sprintf("%s/data/%s", mount, path)

	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{strconv.Itoa(version)}
	}

	secret, err := client.Logical().ReadWithContext(ctx, kvPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %q not found", path)
	}

	rawData, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("secret %q has no data field", path)
	}

	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("secret %q data is not a map", path)
	}

	strData := make(map[string]string, len(dataMap))
	for k, v := range dataMap {
		strData[k] = fmt.Sprintf("%v", v)
	}

	var meta map[string]interface{}
	if rawMeta, ok := secret.Data["metadata"]; ok {
		meta, _ = rawMeta.(map[string]interface{})
	}

	detectedVersion := version
	if meta != nil {
		if v, ok := meta["version"]; ok {
			if vf, ok := v.(float64); ok {
				detectedVersion = int(vf)
			}
		}
	}

	return &SecretVersion{
		Version:  detectedVersion,
		Data:     strData,
		Metadata: meta,
	}, nil
}
