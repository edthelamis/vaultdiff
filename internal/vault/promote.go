package vault

import "fmt"

// PromoteOptions configures a secret promotion between environments.
type PromoteOptions struct {
	DryRun    bool
	Overwrite bool
	Redact    bool
}

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Errors   []string
}

// Summary returns a human-readable summary of the promotion result.
func (r *PromoteResult) Summary() string {
	return fmt.Sprintf("promoted=%d skipped=%d errors=%d",
		len(r.Promoted), len(r.Skipped), len(r.Errors))
}

// PromoteSecrets copies secrets from src to dst, returning a PromoteResult.
// If opts.Overwrite is false, keys already present in dst are skipped.
// If opts.DryRun is true, no writes are performed.
func PromoteSecrets(
	src map[string]string,
	dst map[string]string,
	opts PromoteOptions,
) (*PromoteResult, error) {
	if src == nil {
		return nil, fmt.Errorf("promote: source secrets must not be nil")
	}
	if dst == nil {
		return nil, fmt.Errorf("promote: destination secrets must not be nil")
	}

	result := &PromoteResult{}

	for _, k := range SortedKeys(src) {
		v := src[k]
		if _, exists := dst[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if !opts.DryRun {
			if opts.Redact {
				v = "[REDACTED]"
			}
			dst[k] = v
		}
		result.Promoted = append(result.Promoted, k)
	}

	return result, nil
}
