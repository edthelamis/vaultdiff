package vault

import (
	"strings"
	"testing"
)

func TestCompareSecrets_Matching(t *testing.T) {
	a := map[string]string{"key": "val", "foo": "bar"}
	b := map[string]string{"key": "val", "foo": "bar"}
	res := CompareSecrets(a, b, CompareOptions{})
	if len(res.Matching) != 2 {
		t.Errorf("expected 2 matching, got %d", len(res.Matching))
	}
	if len(res.Differing) != 0 || len(res.OnlyInA) != 0 || len(res.OnlyInB) != 0 {
		t.Error("expected no diffs")
	}
}

func TestCompareSecrets_Differing(t *testing.T) {
	a := map[string]string{"key": "old"}
	b := map[string]string{"key": "new"}
	res := CompareSecrets(a, b, CompareOptions{})
	if len(res.Differing) != 1 || res.Differing[0] != "key" {
		t.Errorf("expected 'key' in differing, got %v", res.Differing)
	}
}

func TestCompareSecrets_OnlyInA(t *testing.T) {
	a := map[string]string{"only_a": "val", "shared": "x"}
	b := map[string]string{"shared": "x"}
	res := CompareSecrets(a, b, CompareOptions{})
	if len(res.OnlyInA) != 1 || res.OnlyInA[0] != "only_a" {
		t.Errorf("expected 'only_a' in OnlyInA, got %v", res.OnlyInA)
	}
}

func TestCompareSecrets_OnlyInB(t *testing.T) {
	a := map[string]string{"shared": "x"}
	b := map[string]string{"shared": "x", "only_b": "val"}
	res := CompareSecrets(a, b, CompareOptions{})
	if len(res.OnlyInB) != 1 || res.OnlyInB[0] != "only_b" {
		t.Errorf("expected 'only_b' in OnlyInB, got %v", res.OnlyInB)
	}
}

func TestCompareSecrets_IgnoreKeys(t *testing.T) {
	a := map[string]string{"key": "old", "ignore": "x"}
	b := map[string]string{"key": "new", "ignore": "y"}
	res := CompareSecrets(a, b, CompareOptions{IgnoreKeys: []string{"ignore"}})
	if len(res.Differing) != 1 || res.Differing[0] != "key" {
		t.Errorf("expected only 'key' differing, got %v", res.Differing)
	}
}

func TestCompareSecrets_OnlyKeys(t *testing.T) {
	a := map[string]string{"key": "val", "other": "old"}
	b := map[string]string{"key": "val", "other": "new"}
	res := CompareSecrets(a, b, CompareOptions{OnlyKeys: []string{"key"}})
	if len(res.Matching) != 1 || res.Matching[0] != "key" {
		t.Errorf("expected only 'key' matching, got %v", res.Matching)
	}
	if len(res.Differing) != 0 {
		t.Errorf("expected no differing keys, got %v", res.Differing)
	}
}

func TestCompareResult_Summary(t *testing.T) {
	res := CompareResult{
		Matching:  []string{"a", "b"},
		Differing: []string{"c"},
		OnlyInA:   []string{"d"},
	}
	summary := res.Summary()
	if !strings.Contains(summary, "Matching") {
		t.Error("summary should contain Matching")
	}
	if !strings.Contains(summary, "Differing") {
		t.Error("summary should contain Differing")
	}
	if !strings.Contains(summary, "Only in A") {
		t.Error("summary should contain Only in A")
	}
}

func TestCompareSecrets_EmptyInputs(t *testing.T) {
	res := CompareSecrets(map[string]string{}, map[string]string{}, CompareOptions{})
	if len(res.Matching)+len(res.Differing)+len(res.OnlyInA)+len(res.OnlyInB) != 0 {
		t.Error("expected empty result for empty inputs")
	}
}
