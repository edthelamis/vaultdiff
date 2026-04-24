package vault

import (
	"strings"
	"testing"
)

func sampleChanges() []DiffEntry {
	return []DiffEntry{
		{Key: "DB_PASS", Status: StatusChanged, OldValue: "old", NewValue: "new"},
		{Key: "API_KEY", Status: StatusAdded, OldValue: "", NewValue: "abc123"},
	}
}

func TestNewAuditLog_Empty(t *testing.T) {
	log := NewAuditLog()
	if len(log.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(log.Entries))
	}
}

func TestAuditLog_Record(t *testing.T) {
	log := NewAuditLog()
	log.Record("production", "secret/app", 1, 2, sampleChanges())

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", e.Environment)
	}
	if e.VersionA != 1 || e.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", e.VersionA, e.VersionB)
	}
	if len(e.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(e.Changes))
	}
}

func TestAuditLog_Summary_NoEntries(t *testing.T) {
	log := NewAuditLog()
	got := log.Summary()
	if got != "No audit entries recorded." {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestAuditLog_Summary_WithEntries(t *testing.T) {
	log := NewAuditLog()
	log.Record("staging", "secret/db", 2, 3, sampleChanges())
	got := log.Summary()
	if !strings.Contains(got, "staging") {
		t.Errorf("summary missing environment name: %q", got)
	}
	if !strings.Contains(got, "secret/db") {
		t.Errorf("summary missing path: %q", got)
	}
	if !strings.Contains(got, "changes=2") {
		t.Errorf("summary missing change count: %q", got)
	}
}

func TestAuditLog_HasChanges_True(t *testing.T) {
	log := NewAuditLog()
	log.Record("dev", "secret/app", 1, 2, sampleChanges())
	if !log.HasChanges() {
		t.Error("expected HasChanges to return true")
	}
}

func TestAuditLog_HasChanges_False(t *testing.T) {
	log := NewAuditLog()
	log.Record("dev", "secret/app", 1, 1, []DiffEntry{})
	if log.HasChanges() {
		t.Error("expected HasChanges to return false")
	}
}
