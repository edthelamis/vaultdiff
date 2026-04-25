package vault

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromoteSecrets_Basic(t *testing.T) {
	src := map[string]string{"DB_PASS": "secret", "API_KEY": "abc"}
	dst := map[string]string{}

	res, err := PromoteSecrets(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if dst["DB_PASS"] != "secret" {
		t.Errorf("expected dst to contain DB_PASS")
	}
}

func TestPromoteSecrets_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}

	res, err := PromoteSecrets(src, dst, PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if dst["KEY"] != "old" {
		t.Errorf("expected dst KEY to remain 'old'")
	}
}

func TestPromoteSecrets_OverwriteReplaces(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}

	_, err := PromoteSecrets(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["KEY"] != "new" {
		t.Errorf("expected dst KEY to be 'new', got %s", dst["KEY"])
	}
}

func TestPromoteSecrets_DryRunNoWrite(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	dst := map[string]string{}

	res, err := PromoteSecrets(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 in promoted list even for dry run")
	}
	if _, ok := dst["KEY"]; ok {
		t.Errorf("expected dst to be empty after dry run")
	}
}

func TestPromoteSecrets_NilSource(t *testing.T) {
	_, err := PromoteSecrets(nil, map[string]string{}, PromoteOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestPromoteReport_Write(t *testing.T) {
	result := &PromoteResult{
		Promoted: []string{"API_KEY"},
		Skipped:  []string{"DB_PASS"},
	}
	report := NewPromoteReport("staging", "production", PromoteOptions{DryRun: false}, result)
	var buf bytes.Buffer
	if err := report.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected source in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected promoted key in output")
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected skipped key in output")
	}
}
