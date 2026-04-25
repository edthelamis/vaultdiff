package vault

import (
	"strings"
	"testing"
)

func TestMergeSecrets_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"a": "1"}
	src := map[string]string{"a": "1", "b": "2"}

	res, err := MergeSecrets(dst, src, MergeOptions{Strategy: MergeStrategySource})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "b" {
		t.Errorf("expected Added=[b], got %v", res.Added)
	}
	if res.Merged["b"] != "2" {
		t.Errorf("expected merged b=2, got %q", res.Merged["b"])
	}
}

func TestMergeSecrets_SourceWins(t *testing.T) {
	dst := map[string]string{"key": "old"}
	src := map[string]string{"key": "new"}

	res, err := MergeSecrets(dst, src, MergeOptions{Strategy: MergeStrategySource})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "key" {
		t.Errorf("expected Overwritten=[key], got %v", res.Overwritten)
	}
	if res.Merged["key"] != "new" {
		t.Errorf("expected merged key=new, got %q", res.Merged["key"])
	}
}

func TestMergeSecrets_DestinationWins(t *testing.T) {
	dst := map[string]string{"key": "old"}
	src := map[string]string{"key": "new"}

	res, err := MergeSecrets(dst, src, MergeOptions{Strategy: MergeStrategyDestination})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "key" {
		t.Errorf("expected Skipped=[key], got %v", res.Skipped)
	}
	if res.Merged["key"] != "old" {
		t.Errorf("expected merged key=old (destination preserved), got %q", res.Merged["key"])
	}
}

func TestMergeSecrets_ErrorOnConflict(t *testing.T) {
	dst := map[string]string{"key": "old"}
	src := map[string]string{"key": "new"}

	res, err := MergeSecrets(dst, src, MergeOptions{Strategy: MergeStrategyError})
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	if !strings.Contains(err.Error(), "key") {
		t.Errorf("expected error to mention conflicting key, got: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMergeSecrets_ExcludeKeys(t *testing.T) {
	dst := map[string]string{}
	src := map[string]string{"a": "1", "secret": "s"}

	res, err := MergeSecrets(dst, src, MergeOptions{
		Strategy:    MergeStrategySource,
		ExcludeKeys: []string{"secret"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Merged["secret"]; ok {
		t.Error("excluded key 'secret' should not appear in merged output")
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "secret" {
		t.Errorf("expected Skipped=[secret], got %v", res.Skipped)
	}
}

func TestMergeSecrets_DryRun(t *testing.T) {
	dst := map[string]string{"a": "1"}
	src := map[string]string{"a": "99", "b": "2"}

	res, err := MergeSecrets(dst, src, MergeOptions{Strategy: MergeStrategySource, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Merged map should still reflect original dst in dry-run mode.
	if res.Merged["a"] != "1" {
		t.Errorf("dry-run: expected a=1 (unchanged), got %q", res.Merged["a"])
	}
	if _, ok := res.Merged["b"]; ok {
		t.Error("dry-run: new key 'b' should not appear in merged map")
	}
	if len(res.Added) != 1 || res.Added[0] != "b" {
		t.Errorf("dry-run: expected Added=[b], got %v", res.Added)
	}
}

func TestMergeResult_Summary(t *testing.T) {
	r := &MergeResult{
		Merged:      map[string]string{"a": "1", "b": "2"},
		Added:       []string{"b"},
		Overwritten: []string{},
		Skipped:     []string{"x"},
		Conflicts:   []string{},
	}
	s := r.Summary()
	if !strings.Contains(s, "merged=2") || !strings.Contains(s, "added=1") {
		t.Errorf("unexpected summary: %q", s)
	}
}
