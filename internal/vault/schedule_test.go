package vault

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestScheduledDiff_DetectsChanges(t *testing.T) {
	calls := 0
	var mu sync.Mutex

	a := map[string]string{"key": "v1"}
	b := map[string]string{"key": "v2"}

	opts := DefaultScheduleOptions()
	opts.Interval = 20 * time.Millisecond
	opts.MaxRuns = 2
	opts.OnDiff = func(entries []DiffEntry) {
		mu.Lock()
		calls++
		mu.Unlock()
	}

	err := ScheduledDiff(context.Background(),
		func() (map[string]string, error) { return a, nil },
		func() (map[string]string, error) { return b, nil },
		opts,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Errorf("expected 2 OnDiff calls, got %d", calls)
	}
}

func TestScheduledDiff_CancelStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	opts := DefaultScheduleOptions()
	opts.Interval = 10 * time.Millisecond
	opts.OnDiff = func(_ []DiffEntry) {}

	go func() {
		time.Sleep(25 * time.Millisecond)
		cancel()
	}()

	err := ScheduledDiff(ctx,
		func() (map[string]string, error) { return map[string]string{"k": "1"}, nil },
		func() (map[string]string, error) { return map[string]string{"k": "2"}, nil },
		opts,
	)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestScheduledDiff_SourceError(t *testing.T) {
	errored := false
	opts := DefaultScheduleOptions()
	opts.Interval = 20 * time.Millisecond
	opts.MaxRuns = 1
	opts.OnError = func(err error) { errored = true }

	_ = ScheduledDiff(context.Background(),
		func() (map[string]string, error) { return nil, errors.New("fetch failed") },
		func() (map[string]string, error) { return map[string]string{}, nil },
		opts,
	)
	if !errored {
		t.Error("expected OnError to be called")
	}
}

func TestScheduledDiff_InvalidInterval(t *testing.T) {
	opts := DefaultScheduleOptions()
	opts.Interval = 0
	err := ScheduledDiff(context.Background(),
		func() (map[string]string, error) { return nil, nil },
		func() (map[string]string, error) { return nil, nil },
		opts,
	)
	if err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestScheduledDiff_NilSource(t *testing.T) {
	opts := DefaultScheduleOptions()
	err := ScheduledDiff(context.Background(), nil, nil, opts)
	if err == nil {
		t.Error("expected error for nil sources")
	}
}
