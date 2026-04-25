package vault

import "testing"

func TestDiffAgainstBaseline_DetectsDrift(t *testing.T) {
	b := NewBaseline("prod", "secret/app", map[string]string{
		"api_key": "old",
		"db_pass": "same",
		"removed": "gone",
	})
	current := map[string]string{
		"api_key": "new",
		"db_pass": "same",
		"added":   "fresh",
	}
	result, err := DiffAgainstBaseline(b, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Environment != "prod" {
		t.Errorf("expected environment prod")
	}
	if result.Path != "secret/app" {
		t.Errorf("expected path secret/app")
	}
	changeMap := map[string]ChangeType{}
	for _, c := range result.Changes {
		changeMap[c.Key] = c.Type
	}
	if changeMap["api_key"] != ChangeModified {
		t.Errorf("expected api_key to be modified")
	}
	if changeMap["removed"] != ChangeRemoved {
		t.Errorf("expected removed to be removed")
	}
	if changeMap["added"] != ChangeAdded {
		t.Errorf("expected added to be added")
	}
	if _, ok := changeMap["db_pass"]; ok {
		t.Errorf("db_pass should not appear in changes")
	}
}

func TestDiffAgainstBaseline_NoDrift(t *testing.T) {
	data := map[string]string{"x": "1", "y": "2"}
	b := NewBaseline("dev", "secret/svc", data)
	result, err := DiffAgainstBaseline(b, map[string]string{"x": "1", "y": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestDiffAgainstBaseline_NilBaseline(t *testing.T) {
	_, err := DiffAgainstBaseline(nil, map[string]string{})
	if err == nil {
		t.Error("expected error for nil baseline")
	}
}

func TestBaselineDiffResult_Summary_WithChanges(t *testing.T) {
	r := &BaselineDiffResult{
		Environment: "prod",
		Path:        "secret/app",
		Changes: []DiffEntry{
			{Key: "a", Type: ChangeAdded},
			{Key: "b", Type: ChangeRemoved},
			{Key: "c", Type: ChangeModified},
		},
	}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
