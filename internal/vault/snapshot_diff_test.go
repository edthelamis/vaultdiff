package vault

import (
	"strings"
	"testing"
)

func makeSnap(env, path string, version int, data map[string]string) *Snapshot {
	return TakeSnapshot(env, path, version, data)
}

func TestDiffSnapshots_DetectsChanges(t *testing.T) {
	from := makeSnap("prod", "secret/app", 1, map[string]string{
		"DB_HOST": "old-host",
		"DB_PORT": "5432",
	})
	to := makeSnap("prod", "secret/app", 2, map[string]string{
		"DB_HOST": "new-host",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
	})

	result, err := DiffSnapshots(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasChanges() {
		t.Error("expected changes to be detected")
	}

	changeMap := make(map[string]ChangeType)
	for _, c := range result.Changes {
		changeMap[c.Key] = c.ChangeType
	}
	if changeMap["DB_HOST"] != Changed {
		t.Errorf("expected DB_HOST to be Changed, got %v", changeMap["DB_HOST"])
	}
	if changeMap["DB_NAME"] != Added {
		t.Errorf("expected DB_NAME to be Added, got %v", changeMap["DB_NAME"])
	}
	if changeMap["DB_PORT"] != Unchanged {
		t.Errorf("expected DB_PORT to be Unchanged, got %v", changeMap["DB_PORT"])
	}
}

func TestDiffSnapshots_NilInput(t *testing.T) {
	_, err := DiffSnapshots(nil, makeSnap("dev", "secret/x", 1, nil))
	if err == nil {
		t.Error("expected error for nil from-snapshot")
	}
	_, err = DiffSnapshots(makeSnap("dev", "secret/x", 1, nil), nil)
	if err == nil {
		t.Error("expected error for nil to-snapshot")
	}
}

func TestDiffSnapshots_NoChanges(t *testing.T) {
	data := map[string]string{"KEY": "val"}
	from := makeSnap("staging", "secret/svc", 4, data)
	to := makeSnap("staging", "secret/svc", 5, data)

	result, err := DiffSnapshots(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestSnapshotDiffResult_Summary(t *testing.T) {
	from := makeSnap("prod", "secret/app", 1, map[string]string{"A": "1", "B": "2"})
	to := makeSnap("prod", "secret/app", 2, map[string]string{"A": "changed", "C": "3"})

	result, _ := DiffSnapshots(from, to)
	summary := result.Summary()

	if !strings.Contains(summary, "v1→v2") {
		t.Errorf("summary missing version range: %q", summary)
	}
	if !strings.Contains(summary, "prod") {
		t.Errorf("summary missing environment: %q", summary)
	}
}
