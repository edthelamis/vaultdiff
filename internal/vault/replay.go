package vault

import (
	"context"
	"fmt"
	"time"
)

// ReplayEvent represents a single replayed diff event from history.
type ReplayEvent struct {
	Timestamp time.Time
	Path      string
	Changes   []DiffEntry
	Version   int
}

// ReplayOptions configures how history replay is performed.
type ReplayOptions struct {
	// Since filters events to only those after this time.
	Since time.Time
	// Until filters events to only those before this time.
	Until time.Time
	// Path restricts replay to a specific secret path.
	Path string
	// MaxEvents caps the number of events emitted (0 = unlimited).
	MaxEvents int
}

// ReplayHistory replays recorded history entries as a stream of ReplayEvents.
// Events are sent to the returned channel in chronological order.
// The caller must consume or cancel via ctx to avoid blocking.
func ReplayHistory(ctx context.Context, h *History, opts ReplayOptions) (<-chan ReplayEvent, error) {
	if h == nil {
		return nil, fmt.Errorf("replay: history must not be nil")
	}

	entries := h.Entries

	// Apply time filters.
	var filtered []HistoryEntry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		if opts.Path != "" && e.Path != opts.Path {
			continue
		}
		filtered = append(filtered, e)
	}

	ch := make(chan ReplayEvent, len(filtered))

	go func() {
		defer close(ch)
		count := 0
		for _, entry := range filtered {
			if opts.MaxEvents > 0 && count >= opts.MaxEvents {
				return
			}
			event := ReplayEvent{
				Timestamp: entry.Timestamp,
				Path:      entry.Path,
				Changes:   entry.Changes,
				Version:   entry.Version,
			}
			select {
			case <-ctx.Done():
				return
			case ch <- event:
				count++
			}
		}
	}()

	return ch, nil
}
