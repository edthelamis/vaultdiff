package vault

import (
	"testing"
)

func baseEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "db/password", Change: Added, NewValue: "secret"},
		{Key: "db/host", Change: Changed, OldValue: "old", NewValue: "new"},
		{Key: "app/token", Change: Removed, OldValue: "tok"},
		{Key: "app/name", Change: Unchanged, OldValue: "myapp", NewValue: "myapp"},
	}
}

func TestEnforcePolicy_NoViolations(t *testing.T) {
	p := Policy{
		Name: "permissive",
		Rules: []PolicyRule{
			{KeyPattern: "*", AllowRead: true, AllowWrite: true, AllowDelete: true},
		},
	}
	v := EnforcePolicy(baseEntries(), p)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestEnforcePolicy_ForbidWrite(t *testing.T) {
	p := Policy{
		Name: "readonly",
		Rules: []PolicyRule{
			{KeyPattern: "db/", AllowRead: true, AllowWrite: false, AllowDelete: true},
		},
	}
	v := EnforcePolicy(baseEntries(), p)
	// db/password (Added) and db/host (Changed) should violate
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestEnforcePolicy_ForbidDelete(t *testing.T) {
	p := Policy{
		Name: "nodelete",
		Rules: []PolicyRule{
			{KeyPattern: "app/", AllowRead: true, AllowWrite: true, AllowDelete: false},
		},
	}
	v := EnforcePolicy(baseEntries(), p)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "app/token" {
		t.Errorf("expected violation on app/token, got %q", v[0].Key)
	}
}

func TestEnforcePolicy_WildcardMatchesAll(t *testing.T) {
	p := Policy{
		Name: "strict",
		Rules: []PolicyRule{
			{KeyPattern: "*", AllowRead: true, AllowWrite: false, AllowDelete: false},
		},
	}
	v := EnforcePolicy(baseEntries(), p)
	// Added, Changed, Removed should all violate
	if len(v) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(v))
	}
}

func TestMatchPattern_Prefix(t *testing.T) {
	if !matchPattern("db/", "db/password") {
		t.Error("expected db/ to match db/password")
	}
	if matchPattern("db/", "app/token") {
		t.Error("expected db/ not to match app/token")
	}
}
