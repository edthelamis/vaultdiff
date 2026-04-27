package vault

import (
	"errors"
	"testing"
)

func TestCloneSecret_Basic(t *testing.T) {
	src := map[string]interface{}{"foo": "bar", "baz": "qux"}
	dest := map[string]interface{}{}
	written := map[string]interface{}{}

	writeFn := func(_ string, data map[string]interface{}) error {
		for k, v := range data {
			written[k] = v
		}
		return nil
	}

	res, err := CloneSecret(src, dest, "secret/src", "secret/dest", writeFn, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.KeysCopied) != 2 {
		t.Errorf("expected 2 keys copied, got %d", len(res.KeysCopied))
	}
	if written["foo"] != "bar" {
		t.Errorf("expected foo=bar in destination")
	}
}

func TestCloneSecret_ExcludeKeys(t *testing.T) {
	src := map[string]interface{}{"foo": "bar", "secret": "topsecret"}
	dest := map[string]interface{}{}
	writeFn := func(_ string, _ map[string]interface{}) error { return nil }

	res, err := CloneSecret(src, dest, "secret/src", "secret/dest", writeFn,
		CloneOptions{ExcludeKeys: []string{"secret"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.KeysCopied) != 1 || res.KeysCopied[0] != "foo" {
		t.Errorf("expected only 'foo' copied, got %v", res.KeysCopied)
	}
	if len(res.KeysSkipped) != 1 || res.KeysSkipped[0] != "secret" {
		t.Errorf("expected 'secret' skipped, got %v", res.KeysSkipped)
	}
}

func TestCloneSecret_NoOverwrite(t *testing.T) {
	src := map[string]interface{}{"foo": "new"}
	dest := map[string]interface{}{"foo": "old"}
	writeFn := func(_ string, _ map[string]interface{}) error { return nil }

	res, err := CloneSecret(src, dest, "secret/src", "secret/dest", writeFn, CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.KeysSkipped) != 1 {
		t.Errorf("expected 1 key skipped due to no-overwrite, got %d", len(res.KeysSkipped))
	}
}

func TestCloneSecret_DryRun(t *testing.T) {
	src := map[string]interface{}{"foo": "bar"}
	dest := map[string]interface{}{}
	called := false
	writeFn := func(_ string, _ map[string]interface{}) error {
		called = true
		return nil
	}

	res, err := CloneSecret(src, dest, "secret/src", "secret/dest", writeFn, CloneOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("writeFn should not be called in dry-run mode")
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
}

func TestCloneSecret_NilSource(t *testing.T) {
	writeFn := func(_ string, _ map[string]interface{}) error { return nil }
	_, err := CloneSecret(nil, map[string]interface{}{}, "src", "dest", writeFn, CloneOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestCloneSecret_WriteError(t *testing.T) {
	src := map[string]interface{}{"foo": "bar"}
	dest := map[string]interface{}{}
	writeFn := func(_ string, _ map[string]interface{}) error {
		return errors.New("write failed")
	}

	_, err := CloneSecret(src, dest, "secret/src", "secret/dest", writeFn, CloneOptions{})
	if err == nil {
		t.Error("expected error from writeFn")
	}
}

func TestCloneResult_Summary(t *testing.T) {
	r := CloneResult{
		SourcePath: "secret/src",
		DestPath:   "secret/dest",
		KeysCopied: []string{"a", "b"},
		KeysSkipped: []string{"c"},
	}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
	r.DryRun = true
	ds := r.Summary()
	if ds == s {
		t.Error("dry-run summary should differ from normal summary")
	}
}
