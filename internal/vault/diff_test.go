package vault

import (
	"testing"
)

func TestDiffSecrets_Added(t *testing.T) {
	old := map[string]string{"key1": "val1"}
	new := map[string]string{"key1": "val1", "key2": "val2"}

	result := DiffSecrets(old, new)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(result.Added))
	}
	if result.Added["key2"] != "val2" {
		t.Errorf("expected added key2=val2, got %s", result.Added["key2"])
	}
}

func TestDiffSecrets_Removed(t *testing.T) {
	old := map[string]string{"key1": "val1", "key2": "val2"}
	new := map[string]string{"key1": "val1"}

	result := DiffSecrets(old, new)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(result.Removed))
	}
	if result.Removed["key2"] != "val2" {
		t.Errorf("expected removed key2=val2, got %s", result.Removed["key2"])
	}
}

func TestDiffSecrets_Changed(t *testing.T) {
	old := map[string]string{"key1": "oldval"}
	new := map[string]string{"key1": "newval"}

	result := DiffSecrets(old, new)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(result.Changed))
	}
	pair := result.Changed["key1"]
	if pair[0] != "oldval" || pair[1] != "newval" {
		t.Errorf("unexpected changed values: %v", pair)
	}
}

func TestDiffSecrets_Unchanged(t *testing.T) {
	old := map[string]string{"key1": "val1"}
	new := map[string]string{"key1": "val1"}

	result := DiffSecrets(old, new)

	if result.HasChanges() {
		t.Error("expected no changes")
	}
	if result.Unchanged["key1"] != "val1" {
		t.Errorf("expected unchanged key1=val1")
	}
}

func TestDiffSecrets_EmptyBoth(t *testing.T) {
	result := DiffSecrets(map[string]string{}, map[string]string{})
	if result.HasChanges() {
		t.Error("expected no changes for empty maps")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]string{"b": "2", "a": "1", "c": "3"}
	keys := SortedKeys(m)
	expected := []string{"a", "b", "c"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, k)
		}
	}
}
