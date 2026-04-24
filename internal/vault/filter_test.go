package vault

import (
	"testing"
)

func sampleEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "db/password", Type: ChangeTypeChanged, OldValue: "old", NewValue: "new"},
		{Key: "db/user", Type: ChangeTypeUnchanged, OldValue: "admin", NewValue: "admin"},
		{Key: "app/secret", Type: ChangeTypeAdded, OldValue: "", NewValue: "abc123"},
		{Key: "app/debug", Type: ChangeTypeRemoved, OldValue: "true", NewValue: ""},
		{Key: "infra/key", Type: ChangeTypeChanged, OldValue: "x", NewValue: "y"},
	}
}

func TestFilterDiff_ByChangeType(t *testing.T) {
	entries := sampleEntries()
	result := FilterDiff(entries, FilterOptions{ChangeTypes: []string{"added", "removed"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Type != ChangeTypeAdded && e.Type != ChangeTypeRemoved {
			t.Errorf("unexpected type %q", e.Type)
		}
	}
}

func TestFilterDiff_ByKeyPrefix(t *testing.T) {
	entries := sampleEntries()
	result := FilterDiff(entries, FilterOptions{KeyPrefix: "db/"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilterDiff_ExcludeKeys(t *testing.T) {
	entries := sampleEntries()
	result := FilterDiff(entries, FilterOptions{ExcludeKeys: []string{"db/password", "infra/key"}})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "db/password" || e.Key == "infra/key" {
			t.Errorf("excluded key %q still present", e.Key)
		}
	}
}

func TestFilterDiff_NoOptions(t *testing.T) {
	entries := sampleEntries()
	result := FilterDiff(entries, FilterOptions{})
	if len(result) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestFilterDiff_CombinedOptions(t *testing.T) {
	entries := sampleEntries()
	result := FilterDiff(entries, FilterOptions{
		ChangeTypes: []string{"changed"},
		KeyPrefix:   "db/",
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "db/password" {
		t.Errorf("unexpected key %q", result[0].Key)
	}
}

func TestFilterDiff_EmptyInput(t *testing.T) {
	result := FilterDiff([]DiffEntry{}, FilterOptions{ChangeTypes: []string{"added"}})
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
