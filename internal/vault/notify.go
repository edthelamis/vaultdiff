package vault

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// NotifyFormat controls how watch notifications are rendered.
type NotifyFormat string

const (
	NotifyFormatText NotifyFormat = "text"
	NotifyFormatJSON NotifyFormat = "json"
)

// NotifyOptions configures notification output.
type NotifyOptions struct {
	Format NotifyFormat
	Writer io.Writer
}

// HandleWatchEvents reads events from ch and writes human-readable or JSON
// notifications to opts.Writer until the channel is closed.
func HandleWatchEvents(ch <-chan WatchEvent, opts NotifyOptions) {
	for event := range ch {
		writeNotification(event, opts)
	}
}

func writeNotification(e WatchEvent, opts NotifyOptions) {
	ts := e.Timestamp.Format(time.RFC3339)

	if e.Err != nil {
		fmt.Fprintf(opts.Writer, "[%s] ERROR env=%s path=%s err=%v\n", ts, e.Environment, e.Path, e.Err)
		return
	}

	switch opts.Format {
	case NotifyFormatJSON:
		writeJSONNotification(e, opts.Writer, ts)
	default:
		writeTextNotification(e, opts.Writer, ts)
	}
}

func writeTextNotification(e WatchEvent, w io.Writer, ts string) {
	fmt.Fprintf(w, "[%s] CHANGE env=%s path=%s version=%d changes=%d\n",
		ts, e.Environment, e.Path, e.Version, len(e.Changes))
	for _, c := range e.Changes {
		fmt.Fprintf(w, "  %s %s\n", strings.ToUpper(string(c.Type)), c.Key)
	}
}

func writeJSONNotification(e WatchEvent, w io.Writer, ts string) {
	keys := make([]string, 0, len(e.Changes))
	for _, c := range e.Changes {
		keys = append(keys, fmt.Sprintf(`{"type":%q,"key":%q}`, c.Type, c.Key))
	}
	fmt.Fprintf(w, `{"timestamp":%q,"env":%q,"path":%q,"version":%d,"changes":[%s]}`+"\n",
		ts, e.Environment, e.Path, e.Version, strings.Join(keys, ","))
}
