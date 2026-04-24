package vault

import (
	"context"
	"time"
)

// WatchOptions configures the polling behavior for secret watching.
type WatchOptions struct {
	Interval  time.Duration
	MaxPolls  int // 0 means unlimited
}

// WatchEvent holds the result of a single poll cycle.
type WatchEvent struct {
	Environment string
	Path        string
	Version     int
	Changes     []DiffEntry
	Err         error
	Timestamp   time.Time
}

// WatchSecret polls a Vault secret path at a given interval, emitting a
// WatchEvent on the returned channel whenever the secret version changes.
func WatchSecret(
	ctx context.Context,
	client *Client,
	env, path string,
	initialVersion int,
	opts WatchOptions,
) <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)

	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	go func() {
		defer close(ch)

		prevVersion := initialVersion
		prevData, _ := GetSecretVersion(client, path, prevVersion)

		polls := 0
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				polls++
				nextVersion := prevVersion + 1
				nextData, err := GetSecretVersion(client, path, nextVersion)
				if err != nil {
					ch <- WatchEvent{Environment: env, Path: path, Err: err, Timestamp: t}
				} else if nextData != nil {
					changes := DiffSecrets(prevData, nextData)
					ch <- WatchEvent{
						Environment: env,
						Path:        path,
						Version:     nextVersion,
						Changes:     changes,
						Timestamp:   t,
					}
					prevVersion = nextVersion
					prevData = nextData
				}
				if opts.MaxPolls > 0 && polls >= opts.MaxPolls {
					return
				}
			}
		}
	}()

	return ch
}
