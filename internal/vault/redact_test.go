package vault

import (
	"testing"
)

func sensitiveEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "db_password", ChangeType: Changed, OldValue: "old123", NewValue: "new456"},
		{Key: "api_token", ChangeType: Added, OldValue: "", NewValue: "tok_abc"},
		{Key: "app_name", ChangeType: Unchanged, OldValue: "myapp", NewValue: "myapp"},
		{Key: "secret_key", ChangeType: Removed, OldValue: "s3cr3t", NewValue: ""},
		{Key: "region", ChangeType: Unchanged, OldValue: "us-east-1", NewValue: "us-east-1"},
	}
}

func TestRedactDiff_SensitiveKeysRedacted(t *testing.T) {
	entries := sensitiveEntries()
	opts := DefaultRedactOptions()
	result := RedactDiff(entries, opts)

	for _, e := range result {
		if isSensitiveKey(e.Key, opts.KeyPatterns) {
			if e.OldValue != "[REDACTED]" || e.NewValue != "[REDACTED]" {
				t.Errorf("key %q: expected redacted values, got old=%q new=%q", e.Key, e.OldValue, e.NewValue)
			}
		}
	}
}

func TestRedactDiff_NonSensitiveKeysUnchanged(t *testing.T) {
	entries := sensitiveEntries()
	opts := DefaultRedactOptions()
	result := RedactDiff(entries, opts)

	for _, e := range result {
		if e.Key == "app_name" && e.OldValue != "myapp" {
			t.Errorf("app_name should not be redacted, got %q", e.OldValue)
		}
		if e.Key == "region" && e.NewValue != "us-east-1" {
			t.Errorf("region should not be redacted, got %q", e.NewValue)
		}
	}
}

func TestRedactDiff_CustomReplacement(t *testing.T) {
	entries := []DiffEntry{
		{Key: "api_token", ChangeType: Changed, OldValue: "old", NewValue: "new"},
	}
	opts := RedactOptions{
		KeyPatterns: []string{"*token*"},
		Replacement: "***",
	}
	result := RedactDiff(entries, opts)
	if result[0].OldValue != "***" || result[0].NewValue != "***" {
		t.Errorf("expected '***' replacement, got old=%q new=%q", result[0].OldValue, result[0].NewValue)
	}
}

func TestRedactDiff_EmptyPatterns(t *testing.T) {
	entries := sensitiveEntries()
	opts := RedactOptions{KeyPatterns: []string{}, Replacement: "[REDACTED]"}
	result := RedactDiff(entries, opts)
	for i, e := range result {
		if e.OldValue != entries[i].OldValue || e.NewValue != entries[i].NewValue {
			t.Errorf("key %q: values should be unchanged with no patterns", e.Key)
		}
	}
}

func TestGlobMatch(t *testing.T) {
	cases := []struct {
		pattern string
		s       string
		want    bool
	}{
		{"*password*", "db_password", true},
		{"*token*", "api_token", true},
		{"*secret*", "app_name", false},
		{"exact", "exact", true},
		{"exact", "not_exact", false},
	}
	for _, tc := range cases {
		got := globMatch(tc.pattern, tc.s)
		if got != tc.want {
			t.Errorf("globMatch(%q, %q) = %v, want %v", tc.pattern, tc.s, got, tc.want)
		}
	}
}
