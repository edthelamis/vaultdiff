package vault

import (
	"fmt"
	"strings"
)

// PolicyRule defines access rules for a secret key pattern.
type PolicyRule struct {
	KeyPattern string
	AllowRead  bool
	AllowWrite bool
	AllowDelete bool
}

// Policy holds a named set of rules governing secret access.
type Policy struct {
	Name  string
	Rules []PolicyRule
}

// PolicyViolation describes a diff entry that violates a policy.
type PolicyViolation struct {
	Key    string
	Change ChangeType
	Reason string
}

// EnforcePolicy checks a slice of DiffEntries against a Policy and returns
// any violations found.
func EnforcePolicy(entries []DiffEntry, p Policy) []PolicyViolation {
	var violations []PolicyViolation
	for _, entry := range entries {
		for _, rule := range p.Rules {
			if !matchPattern(rule.KeyPattern, entry.Key) {
				continue
			}
			switch entry.Change {
			case Added:
				if !rule.AllowWrite {
					violations = append(violations, PolicyViolation{
						Key:    entry.Key,
						Change: entry.Change,
						Reason: fmt.Sprintf("policy %q forbids write on key %q", p.Name, entry.Key),
					})
				}
			case Removed:
				if !rule.AllowDelete {
					violations = append(violations, PolicyViolation{
						Key:    entry.Key,
						Change: entry.Change,
						Reason: fmt.Sprintf("policy %q forbids delete on key %q", p.Name, entry.Key),
					})
				}
			case Changed:
				if !rule.AllowWrite {
					violations = append(violations, PolicyViolation{
						Key:    entry.Key,
						Change: entry.Change,
						Reason: fmt.Sprintf("policy %q forbids write on key %q", p.Name, entry.Key),
					})
				}
			}
		}
	}
	return violations
}

// matchPattern returns true if key starts with the given pattern prefix.
// A pattern of "*" matches all keys.
func matchPattern(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	return strings.HasPrefix(key, pattern)
}
