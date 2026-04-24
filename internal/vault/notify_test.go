package vault

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleEvent(errMsg string) WatchEvent {
	e := WatchEvent{
		Environment: "staging",
		Path:        "secret/data/app",
		Version:     3,
		Timestamp:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Changes: []DiffEntry{
			{Key: "DB_PASS", Type: ChangeTypeChanged, OldValue: "old", NewValue: "new"},
			{Key: "NEW_KEY", Type: ChangeTypeAdded, NewValue: "val"},
		},
	}
	if errMsg != "" {
		e.Err = &watchError{msg: errMsg}
	}
	return e
}

type watchError struct{ msg string }

func (e *watchError) Error() string { return e.msg }

func TestHandleWatchEvents_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	ch := make(chan WatchEvent, 1)
	ch <- sampleEvent("")
	close(ch)

	HandleWatchEvents(ch, NotifyOptions{Format: NotifyFormatText, Writer: &buf})

	out := buf.String()
	if !strings.Contains(out, "CHANGE") {
		t.Errorf("expected CHANGE in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
}

func TestHandleWatchEvents_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	ch := make(chan WatchEvent, 1)
	ch <- sampleEvent("")
	close(ch)

	HandleWatchEvents(ch, NotifyOptions{Format: NotifyFormatJSON, Writer: &buf})

	out := buf.String()
	if !strings.Contains(out, `"env":"staging"`) {
		t.Errorf("expected env field, got: %s", out)
	}
	if !strings.Contains(out, `"version":3`) {
		t.Errorf("expected version field, got: %s", out)
	}
}

func TestHandleWatchEvents_ErrorEvent(t *testing.T) {
	var buf bytes.Buffer
	ch := make(chan WatchEvent, 1)
	ch <- sampleEvent("connection refused")
	close(ch)

	HandleWatchEvents(ch, NotifyOptions{Format: NotifyFormatText, Writer: &buf})

	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", buf.String())
	}
}

func TestHandleWatchEvents_EmptyChannel(t *testing.T) {
	var buf bytes.Buffer
	ch := make(chan WatchEvent)
	close(ch)
	HandleWatchEvents(ch, NotifyOptions{Format: NotifyFormatText, Writer: &buf})
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty channel")
	}
}
