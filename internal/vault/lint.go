package vault

import (
	"fmt"
	"strings"
)

// LintRule defines a rule that validates secret keys or values.
type LintRule struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Message string `json:"message"`
	Target  string `json:"target"` // "key" or "value"
}

// LintViolation represents a single rule violation found during linting.
type LintViolation struct {
	Rule    string
	Key     string
	Message string
}

// LintResult holds all violations found for a given path.
type LintResult struct {
	Path       string
	Violations []LintViolation
}

// HasViolations returns true if any violations were found.
func (r *LintResult) HasViolations() bool {
	return len(r.Violations) > 0
}

// Summary returns a human-readable summary of the lint result.
func (r *LintResult) Summary() string {
	if !r.HasViolations() {
		return fmt.Sprintf("[lint] %s: OK", r.Path)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[lint] %s: %d violation(s)\n", r.Path, len(r.Violations)))
	for _, v := range r.Violations {
		sb.WriteString(fmt.Sprintf("  - [%s] key=%q: %s\n", v.Rule, v.Key, v.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// LintSecrets applies a set of lint rules to a map of secret key/value pairs.
func LintSecrets(path string, secrets map[string]string, rules []LintRule) LintResult {
	result := LintResult{Path: path}
	for _, rule := range rules {
		for k, v := range secrets {
			var subject string
			if rule.Target == "value" {
				subject = v
			} else {
				subject = k
			}
			if globMatch(rule.Pattern, subject) {
				result.Violations = append(result.Violations, LintViolation{
					Rule:    rule.Name,
					Key:     k,
					Message: rule.Message,
				})
			}
		}
	}
	return result
}
