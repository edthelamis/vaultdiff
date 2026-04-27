package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLockIndex_AcquireAndIsLocked(t *testing.T) {
	li := NewLockIndex()
	if err := li.Acquire("secret/db", "alice", "deploying", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !li.IsLocked("secret/db") {
		t.Error("expected path to be locked")
	}
}

func TestLockIndex_AcquireDuplicate(t *testing.T) {
	li := NewLockIndex()
	_ = li.Acquire("secret/db", "alice", "reason", 0)
	err := li.Acquire("secret/db", "bob", "other", 0)
	if err == nil {
		t.Fatal("expected error on duplicate lock")
	}
}

func TestLockIndex_AcquireExpired(t *testing.T) {
	li := NewLockIndex()
	// Acquire with a TTL already in the past by manipulating the entry directly.
	li.Locks["secret/db"] = LockEntry{
		Path:      "secret/db",
		Owner:     "alice",
		LockedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := li.Acquire("secret/db", "bob", "fresh lock", 0); err != nil {
		t.Fatalf("expected expired lock to be replaceable: %v", err)
	}
	if li.Locks["secret/db"].Owner != "bob" {
		t.Error("expected owner to be updated to bob")
	}
}

func TestLockIndex_Release(t *testing.T) {
	li := NewLockIndex()
	_ = li.Acquire("secret/api", "alice", "", 0)
	if err := li.Release("secret/api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if li.IsLocked("secret/api") {
		t.Error("expected path to be unlocked after release")
	}
}

func TestLockIndex_ReleaseNotLocked(t *testing.T) {
	li := NewLockIndex()
	err := li.Release("secret/missing")
	if err == nil {
		t.Fatal("expected error when releasing non-existent lock")
	}
}

func TestLockIndex_IsLocked_WithActiveTTL(t *testing.T) {
	li := NewLockIndex()
	_ = li.Acquire("secret/ttl", "alice", "", 10*time.Minute)
	if !li.IsLocked("secret/ttl") {
		t.Error("expected path to be locked within TTL")
	}
}

func TestLockIndex_IsLocked_ExpiredTTL(t *testing.T) {
	li := NewLockIndex()
	// Insert an entry whose TTL has already elapsed.
	li.Locks["secret/expired"] = LockEntry{
		Path:      "secret/expired",
		Owner:     "alice",
		LockedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if li.IsLocked("secret/expired") {
		t.Error("expected expired lock to be treated as unlocked")
	}
}

func TestSaveAndLoadLockIndex_RoundTrip(t *testing.T) {
	li := NewLockIndex()
	_ = li.Acquire("secret/x", "carol", "testing", 0)

	tmp := filepath.Join(t.TempDir(), "locks.json")
	if err := SaveLockIndex(li, tmp); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadLockIndex(tmp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !loaded.IsLocked("secret/x") {
		t.Error("expected lock to survive round-trip")
	}
	if loaded.Locks["secret/x"].Owner != "carol" {
		t.Errorf("expected owner carol, got %s", loaded.Locks["secret/x"].Owner)
	}
}

func TestLoadLockIndex_MissingFile(t *testing.T) {
	_, err := LoadLockIndex(filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveLockIndex_InvalidPath(t *testing.T) {
	li := NewLockIndex()
	err := SaveLockIndex(li, "/nonexistent/dir/locks.json")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
	_ = os.Remove("/nonexistent/dir/locks.json")
}
