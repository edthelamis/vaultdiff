package vault

import (
	"context"
	"fmt"
	"time"
)

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Path        string
	FromVersion int
	ToVersion   int
	Keys        []string
	RolledBack  bool
	DryRun      bool
	Timestamp   time.Time
}

// RollbackOptions configures rollback behaviour.
type RollbackOptions struct {
	DryRun    bool
	MountPath string
}

// DefaultRollbackOptions returns sensible defaults.
func DefaultRollbackOptions() RollbackOptions {
	return RollbackOptions{
		MountPath: "secret",
	}
}

// RollbackSecret reverts a KV-v2 secret at path to the given target version.
// If DryRun is true the write is skipped and the result is marked accordingly.
func RollbackSecret(ctx context.Context, client *Client, path string, targetVersion int, opts RollbackOptions) (*RollbackResult, error) {
	if client == nil {
		return nil, fmt.Errorf("rollback: client must not be nil")
	}
	if path == "" {
		return nil, fmt.Errorf("rollback: path must not be empty")
	}
	if targetVersion < 1 {
		return nil, fmt.Errorf("rollback: targetVersion must be >= 1, got %d", targetVersion)
	}

	// Read the target version to obtain its data.
	data, currentVersion, err := readVersionedSecret(ctx, client, opts.MountPath, path, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("rollback: reading version %d of %q: %w", targetVersion, path, err)
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	result := &RollbackResult{
		Path:        path,
		FromVersion: currentVersion,
		ToVersion:   targetVersion,
		Keys:        keys,
		DryRun:      opts.DryRun,
		Timestamp:   time.Now().UTC(),
	}

	if opts.DryRun {
		result.RolledBack = false
		return result, nil
	}

	if err := writeSecret(ctx, client, opts.MountPath, path, data); err != nil {
		return nil, fmt.Errorf("rollback: writing restored data to %q: %w", path, err)
	}

	result.RolledBack = true
	return result, nil
}

// readVersionedSecret fetches data and the current (latest) version number for a KV-v2 path.
func readVersionedSecret(ctx context.Context, c *Client, mount, path string, version int) (map[string]string, int, error) {
	secret, err := GetSecretVersion(ctx, c, mount, path, version)
	if err != nil {
		return nil, 0, err
	}
	current, _ := GetSecretVersion(ctx, c, mount, path, 0)
	currentVer := version
	if current != nil {
		if v, ok := current["__version"]; ok {
			fmt.Sscanf(v, "%d", &currentVer)
		}
	}
	delete(secret, "__version")
	return secret, currentVer, nil
}

// writeSecret writes data back to a KV-v2 path.
func writeSecret(ctx context.Context, c *Client, mount, path string, data map[string]string) error {
	writePath := fmt.Sprintf("%s/data/%s", mount, path)
	payload := map[string]interface{}{"data": toInterfaceMap(data)}
	_, err := c.Logical().WriteWithContext(ctx, writePath, payload)
	return err
}

func toInterfaceMap(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
