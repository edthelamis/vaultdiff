package vault

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

var alertEntries = []DiffEntry{
	{Key: "db/password", ChangeType: "changed", OldValue: "old", NewValue: "new"},
	{Key: "app/debug", ChangeType: "added", OldValue: "", NewValue: "true"},
	{Key: "app/name", ChangeType: "unchanged", OldValue: "myapp", NewValue: "myapp"},
	{Key: "db/host", ChangeType: "removed", OldValue: "localhost", NewValue: ""},
}

func TestEvaluateAlerts_MatchesChangeType(t *testing.T) {
	rules := []AlertRule{
		{Name: "detect-removal", KeyPattern: "*", ChangeType: "removed", Severity: SeverityCritical},
	}
	alerts := EvaluateAlerts(alertEntries, rules)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Key != "db/host" {
		t.Errorf("expected key db/host, got %s", alerts[0].Key)
	}
	if alerts[0].Severity != SeverityCritical {
		t.Errorf("expected CRITICAL severity")
	}
}

func TestEvaluateAlerts_MatchesKeyPattern(t *testing.T) {
	rules := []AlertRule{
		{Name: "db-changes", KeyPattern: "db/*", ChangeType: "", Severity: SeverityWarning},
	}
	alerts := EvaluateAlerts(alertEntries, rules)
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts for db/* pattern, got %d", len(alerts))
	}
}

func TestEvaluateAlerts_NoMatch(t *testing.T) {
	rules := []AlertRule{
		{Name: "no-match", KeyPattern: "nonexistent/*", ChangeType: "added", Severity: SeverityInfo},
	}
	alerts := EvaluateAlerts(alertEntries, rules)
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestEvaluateAlerts_EmptyRules(t *testing.T) {
	alerts := EvaluateAlerts(alertEntries, nil)
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts with no rules")
	}
}

func TestWriteAlerts_TextFormat(t *testing.T) {
	rules := []AlertRule{
		{Name: "any-change", KeyPattern: "*", ChangeType: "changed", Severity: SeverityWarning},
	}
	alerts := EvaluateAlerts(alertEntries, rules)
	var buf bytes.Buffer
	if err := WriteAlerts(&buf, alerts, "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[WARNING]") {
		t.Errorf("expected [WARNING] in text output, got: %s", out)
	}
	if !strings.Contains(out, "db/password") {
		t.Errorf("expected key db/password in output")
	}
}

func TestWriteAlerts_JSONFormat(t *testing.T) {
	rules := []AlertRule{
		{Name: "added-keys", KeyPattern: "app/*", ChangeType: "added", Severity: SeverityInfo},
	}
	alerts := EvaluateAlerts(alertEntries, rules)
	var buf bytes.Buffer
	if err := WriteAlerts(&buf, alerts, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []Alert
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded) != 1 || decoded[0].Key != "app/debug" {
		t.Errorf("unexpected decoded alerts: %+v", decoded)
	}
}
