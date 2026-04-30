package vault

import (
	"strings"
	"testing"
)

var lintRules = []LintRule{
	{Name: "no-empty-value", Pattern: "", Target: "value", Message: "value must not be empty"},
	{Name: "no-debug-key", Pattern: "debug_*", Target: "key", Message: "debug keys are not allowed in production"},
}

func TestLintSecrets_NoViolations(t *testing.T) {
	secrets := map[string]string{"api_key": "abc123", "db_pass": "secret"}
	result := LintSecrets("secret/prod", secrets, lintRules)
	if result.HasViolations() {
		t.Errorf("expected no violations, got %d", len(result.Violations))
	}
}

func TestLintSecrets_EmptyValueViolation(t *testing.T) {
	secrets := map[string]string{"api_key": ""}
	result := LintSecrets("secret/prod", secrets, lintRules)
	if !result.HasViolations() {
		t.Fatal("expected violation for empty value")
	}
	if result.Violations[0].Rule != "no-empty-value" {
		t.Errorf("unexpected rule: %s", result.Violations[0].Rule)
	}
}

func TestLintSecrets_KeyPatternViolation(t *testing.T) {
	secrets := map[string]string{"debug_mode": "true"}
	result := LintSecrets("secret/prod", secrets, lintRules)
	if !result.HasViolations() {
		t.Fatal("expected violation for debug key")
	}
	if result.Violations[0].Key != "debug_mode" {
		t.Errorf("unexpected key: %s", result.Violations[0].Key)
	}
}

func TestLintResult_Summary_NoViolations(t *testing.T) {
	result := LintResult{Path: "secret/dev"}
	summary := result.Summary()
	if !strings.Contains(summary, "OK") {
		t.Errorf("expected OK in summary, got: %s", summary)
	}
}

func TestLintResult_Summary_WithViolations(t *testing.T) {
	result := LintResult{
		Path: "secret/prod",
		Violations: []LintViolation{
			{Rule: "no-debug-key", Key: "debug_mode", Message: "not allowed"},
		},
	}
	summary := result.Summary()
	if !strings.Contains(summary, "1 violation") {
		t.Errorf("expected violation count in summary, got: %s", summary)
	}
}

func TestLintSecrets_EmptyRules(t *testing.T) {
	secrets := map[string]string{"key": "val"}
	result := LintSecrets("secret/any", secrets, nil)
	if result.HasViolations() {
		t.Error("expected no violations with empty rules")
	}
}
