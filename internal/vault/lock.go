package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LockEntry represents a lock held on a secret path.
type LockEntry struct {
	Path      string    `json:"path"`
	Owner     string    `json:"owner"`
	Reason    string    `json:"reason"`
	LockedAt  time.Time `json:"locked_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// LockIndex holds all active locks keyed by secret path.
type LockIndex struct {
	Locks map[string]LockEntry `json:"locks"`
}

// NewLockIndex creates an empty LockIndex.
func NewLockIndex() *LockIndex {
	return &LockIndex{Locks: make(map[string]LockEntry)}
}

// Acquire adds a lock for the given path. Returns an error if already locked.
func (li *LockIndex) Acquire(path, owner, reason string, ttl time.Duration) error {
	if entry, exists := li.Locks[path]; exists {
		if entry.ExpiresAt.IsZero() || time.Now().Before(entry.ExpiresAt) {
			return fmt.Errorf("path %q is already locked by %q since %s", path, entry.Owner, entry.LockedAt.Format(time.RFC3339))
		}
	}
	entry := LockEntry{
		Path:     path,
		Owner:    owner,
		Reason:   reason,
		LockedAt: time.Now().UTC(),
	}
	if ttl > 0 {
		entry.ExpiresAt = entry.LockedAt.Add(ttl)
	}
	li.Locks[path] = entry
	return nil
}

// Release removes the lock for the given path. Returns an error if not locked.
func (li *LockIndex) Release(path string) error {
	if _, exists := li.Locks[path]; !exists {
		return fmt.Errorf("path %q is not locked", path)
	}
	delete(li.Locks, path)
	return nil
}

// IsLocked reports whether a path currently has an active lock.
func (li *LockIndex) IsLocked(path string) bool {
	entry, exists := li.Locks[path]
	if !exists {
		return false
	}
	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		return false
	}
	return true
}

// SaveLockIndex persists the lock index to a JSON file.
func SaveLockIndex(li *LockIndex, path string) error {
	data, err := json.MarshalIndent(li, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lock index: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadLockIndex reads a lock index from a JSON file.
func LoadLockIndex(path string) (*LockIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lock index: %w", err)
	}
	var li LockIndex
	if err := json.Unmarshal(data, &li); err != nil {
		return nil, fmt.Errorf("unmarshal lock index: %w", err)
	}
	if li.Locks == nil {
		li.Locks = make(map[string]LockEntry)
	}
	return &li, nil
}
