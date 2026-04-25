package vault

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewCompareReport_Fields(t *testing.T) {
	res := CompareResult{Matching: []string{"a"}}
	rep := NewCompareReport("staging", "prod", res)
	if rep.EnvA != "staging" || rep.EnvB != "prod" {
		t.Errorf("unexpected env names: %s, %s", rep.EnvA, rep.EnvB)
	}
	if len(rep.Result.Matching) != 1 {
		t.Error("expected result to be stored")
	}
}

func TestCompareReport_Write_Header(t *testing.T) {
	res := CompareResult{}
	rep := NewCompareReport("dev", "prod", res)
	var buf bytes.Buffer
	if err := rep.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "dev vs prod") {
		t.Errorf("expected header with env names, got: %s", buf.String())
	}
}

func TestCompareReport_Write_Matching(t *testing.T) {
	res := CompareResult{Matching: []string{"db_pass", "api_key"}}
	rep := NewCompareReport("a", "b", res)
	var buf bytes.Buffer
	_ = rep.Write(&buf)
	out := buf.String()
	if !strings.Contains(out, "MATCHING") {
		t.Error("expected MATCHING section")
	}
	if !strings.Contains(out, "db_pass") || !strings.Contains(out, "api_key") {
		t.Error("expected matching keys in output")
	}
}

func TestCompareReport_Write_Differing(t *testing.T) {
	res := CompareResult{Differing: []string{"secret_key"}}
	rep := NewCompareReport("a", "b", res)
	var buf bytes.Buffer
	_ = rep.Write(&buf)
	out := buf.String()
	if !strings.Contains(out, "DIFFERING") {
		t.Error("expected DIFFERING section")
	}
	if !strings.Contains(out, "~ secret_key") {
		t.Error("expected differing key with ~ prefix")
	}
}

func TestCompareReport_Write_OnlyInA(t *testing.T) {
	res := CompareResult{OnlyInA: []string{"legacy_key"}}
	rep := NewCompareReport("staging", "prod", res)
	var buf bytes.Buffer
	_ = rep.Write(&buf)
	out := buf.String()
	if !strings.Contains(out, "ONLY IN STAGING") {
		t.Error("expected ONLY IN STAGING section")
	}
	if !strings.Contains(out, "- legacy_key") {
		t.Error("expected key with - prefix")
	}
}

func TestCompareReport_Write_OnlyInB(t *testing.T) {
	res := CompareResult{OnlyInB: []string{"new_feature_key"}}
	rep := NewCompareReport("staging", "prod", res)
	var buf bytes.Buffer
	_ = rep.Write(&buf)
	out := buf.String()
	if !strings.Contains(out, "ONLY IN PROD") {
		t.Error("expected ONLY IN PROD section")
	}
	if !strings.Contains(out, "+ new_feature_key") {
		t.Error("expected key with + prefix")
	}
}
