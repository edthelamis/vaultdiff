package vault

import (
	"context"
	"testing"
	"time"
)

func buildReplayHistory() *History {
	h := NewHistory()
	now := time.Now()
	h.Entries = []HistoryEntry{
		{
			Timestamp: now.Add(-3 * time.Hour),
			Path:      "secret/app",
			Version:   1,
			Changes:   []DiffEntry{{Key: "foo", ChangeType: "added"}},
		},
		{
			Timestamp: now.Add(-2 * time.Hour),
			Path:      "secret/app",
			Version:   2,
			Changes:   []DiffEntry{{Key: "bar", ChangeType: "changed"}},
		},
		{
			Timestamp: now.Add(-1 * time.Hour),
			Path:      "secret/other",
			Version:   1,
			Changes:   []DiffEntry{{Key: "baz", ChangeType: "removed"}},
		},
	}
	return h
}

func collectEvents(ch <-chan ReplayEvent) []ReplayEvent {
	var events []ReplayEvent
	for e := range ch {
		events = append(events, e)
	}
	return events
}

func TestReplayHistory_AllEvents(t *testing.T) {
	h := buildReplayHistory()
	ch, err := ReplayHistory(context.Background(), h, ReplayOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := collectEvents(ch)
	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}

func TestReplayHistory_FilterByPath(t *testing.T) {
	h := buildReplayHistory()
	ch, err := ReplayHistory(context.Background(), h, ReplayOptions{Path: "secret/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := collectEvents(ch)
	if len(events) != 2 {
		t.Errorf("expected 2 events for secret/app, got %d", len(events))
	}
}

func TestReplayHistory_FilterBySince(t *testing.T) {
	h := buildReplayHistory()
	opts := ReplayOptions{Since: time.Now().Add(-90 * time.Minute)}
	ch, err := ReplayHistory(context.Background(), h, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := collectEvents(ch)
	if len(events) != 1 {
		t.Errorf("expected 1 event after since filter, got %d", len(events))
	}
}

func TestReplayHistory_MaxEvents(t *testing.T) {
	h := buildReplayHistory()
	ch, err := ReplayHistory(context.Background(), h, ReplayOptions{MaxEvents: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := collectEvents(ch)
	if len(events) != 2 {
		t.Errorf("expected 2 events due to MaxEvents cap, got %d", len(events))
	}
}

func TestReplayHistory_NilHistory(t *testing.T) {
	_, err := ReplayHistory(context.Background(), nil, ReplayOptions{})
	if err == nil {
		t.Error("expected error for nil history, got nil")
	}
}

func TestReplayHistory_CancelContext(t *testing.T) {
	h := buildReplayHistory()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch, err := ReplayHistory(ctx, h, ReplayOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should drain without panic; may have 0 or more events depending on scheduling.
	_ = collectEvents(ch)
}
