package vault

import (
	"context"
	"fmt"
	"time"
)

// ScheduleOptions configures a scheduled diff run.
type ScheduleOptions struct {
	Interval  time.Duration
	MaxRuns   int // 0 means run indefinitely
	OnDiff    func(entries []DiffEntry)
	OnError   func(err error)
}

// DefaultScheduleOptions returns sensible defaults for a scheduled diff.
func DefaultScheduleOptions() ScheduleOptions {
	return ScheduleOptions{
		Interval: 5 * time.Minute,
		MaxRuns:  0,
		OnDiff:   func(entries []DiffEntry) {},
		OnError:  func(err error) {},
	}
}

// ScheduledDiff runs DiffSecrets between two secret maps on a fixed interval.
// It calls opts.OnDiff whenever changes are detected and opts.OnError on failure.
// The loop exits when ctx is cancelled or MaxRuns is reached.
func ScheduledDiff(
	ctx context.Context,
	sourceA, sourceB func() (map[string]string, error),
	opts ScheduleOptions,
) error {
	if opts.Interval <= 0 {
		return fmt.Errorf("schedule interval must be positive")
	}
	if sourceA == nil || sourceB == nil {
		return fmt.Errorf("sourceA and sourceB must not be nil")
	}

	runs := 0
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	runOnce := func() {
		a, err := sourceA()
		if err != nil {
			opts.OnError(fmt.Errorf("sourceA: %w", err))
			return
		}
		b, err := sourceB()
		if err != nil {
			opts.OnError(fmt.Errorf("sourceB: %w", err))
			return
		}
		entries := DiffSecrets(a, b)
		var changed []DiffEntry
		for _, e := range entries {
			if e.ChangeType != "unchanged" {
				changed = append(changed, e)
			}
		}
		if len(changed) > 0 {
			opts.OnDiff(changed)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			runOnce()
			runs++
			if opts.MaxRuns > 0 && runs >= opts.MaxRuns {
				return nil
			}
		}
	}
}
