package vault

import (
	"regexp"
	"strings"
)

// RedactOptions controls which keys are redacted in diff output.
type RedactOptions struct {
	// KeyPatterns is a list of glob-style patterns (e.g. "*password*", "secret_*").
	KeyPatterns []string
	// Replacement is the string used in place of redacted values.
	Replacement string
}

// DefaultRedactOptions returns sensible defaults for redaction.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		KeyPatterns: []string{"*password*", "*secret*", "*token*", "*key*"},
		Replacement: "[REDACTED]",
	}
}

// RedactDiff applies redaction rules to a slice of DiffEntry values,
// replacing sensitive values with the configured replacement string.
func RedactDiff(entries []DiffEntry, opts RedactOptions) []DiffEntry {
	if len(opts.Replacement) == 0 {
		opts.Replacement = "[REDACTED]"
	}
	result := make([]DiffEntry, len(entries))
	for i, e := range entries {
		if isSensitiveKey(e.Key, opts.KeyPatterns) {
			e.OldValue = opts.Replacement
			e.NewValue = opts.Replacement
		}
		result[i] = e
	}
	return result
}

// isSensitiveKey returns true if the key matches any of the given glob patterns.
func isSensitiveKey(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if globMatch(strings.ToLower(p), lower) {
			return true
		}
	}
	return false
}

// globMatch performs simple wildcard matching where '*' matches any substring.
func globMatch(pattern, s string) bool {
	// Convert glob pattern to regex.
	regexStr := "^" + regexp.QuoteMeta(pattern) + "$"
	regexStr = strings.ReplaceAll(regexStr, `\*`, `.*`)
	matched, err := regexp.MatchString(regexStr, s)
	if err != nil {
		return false
	}
	return matched
}
