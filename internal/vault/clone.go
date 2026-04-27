package vault

import (
	"errors"
	"fmt"
)

// CloneOptions configures the behaviour of CloneSecret.
type CloneOptions struct {
	// Overwrite controls whether an existing destination path is replaced.
	Overwrite bool
	// ExcludeKeys lists keys that must not be copied to the destination.
	ExcludeKeys []string
	// DryRun reports what would happen without writing anything.
	DryRun bool
}

// CloneResult describes the outcome of a single clone operation.
type CloneResult struct {
	SourcePath string
	DestPath   string
	KeysCopied []string
	KeysSkipped []string
	DryRun     bool
}

// Summary returns a human-readable one-liner for the clone result.
func (r CloneResult) Summary() string {
	if r.DryRun {
		return fmt.Sprintf("[dry-run] would clone %d key(s) from %s → %s (%d skipped)",
			len(r.KeysCopied), r.SourcePath, r.DestPath, len(r.KeysSkipped))
	}
	return fmt.Sprintf("cloned %d key(s) from %s → %s (%d skipped)",
		len(r.KeysCopied), r.SourcePath, r.DestPath, len(r.KeysSkipped))
}

// CloneSecret copies secrets from srcPath to destPath using the provided
// read/write functions, respecting the supplied CloneOptions.
func CloneSecret(
	src map[string]interface{},
	dest map[string]interface{},
	srcPath, destPath string,
	writeFn func(path string, data map[string]interface{}) error,
	opts CloneOptions,
) (CloneResult, error) {
	if src == nil {
		return CloneResult{}, errors.New("source secret data is nil")
	}
	if srcPath == "" {
		return CloneResult{}, errors.New("source path must not be empty")
	}
	if destPath == "" {
		return CloneResult{}, errors.New("destination path must not be empty")
	}

	excludeSet := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = struct{}{}
	}

	merged := make(map[string]interface{})
	if opts.Overwrite {
		for k, v := range dest {
			merged[k] = v
		}
	} else if dest != nil {
		for k, v := range dest {
			merged[k] = v
		}
	}

	var copied, skipped []string
	for k, v := range src {
		if _, excluded := excludeSet[k]; excluded {
			skipped = append(skipped, k)
			continue
		}
		if _, exists := dest[k]; exists && !opts.Overwrite {
			skipped = append(skipped, k)
			continue
		}
		merged[k] = v
		copied = append(copied, k)
	}

	result := CloneResult{
		SourcePath:  srcPath,
		DestPath:    destPath,
		KeysCopied:  copied,
		KeysSkipped: skipped,
		DryRun:      opts.DryRun,
	}

	if !opts.DryRun && len(copied) > 0 {
		if err := writeFn(destPath, merged); err != nil {
			return CloneResult{}, fmt.Errorf("write to %s failed: %w", destPath, err)
		}
	}

	return result, nil
}
