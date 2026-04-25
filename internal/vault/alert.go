package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// AlertSeverity represents the urgency level of an alert.
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "INFO"
	SeverityWarning  AlertSeverity = "WARNING"
	SeverityCritical AlertSeverity = "CRITICAL"
)

// AlertRule defines criteria for triggering an alert on a diff entry.
type AlertRule struct {
	Name       string        `json:"name"`
	KeyPattern string        `json:"key_pattern"`
	ChangeType string        `json:"change_type"` // added, removed, changed, or "" for any
	Severity   AlertSeverity `json:"severity"`
}

// Alert is a triggered alert produced by evaluating rules against diff entries.
type Alert struct {
	Rule      string        `json:"rule"`
	Severity  AlertSeverity `json:"severity"`
	Key       string        `json:"key"`
	ChangeType string       `json:"change_type"`
	TriggeredAt time.Time   `json:"triggered_at"`
}

// EvaluateAlerts checks a slice of DiffEntry values against the provided rules
// and returns any alerts that match.
func EvaluateAlerts(entries []DiffEntry, rules []AlertRule) []Alert {
	var alerts []Alert
	now := time.Now().UTC()

	for _, entry := range entries {
		for _, rule := range rules {
			if rule.ChangeType != "" && rule.ChangeType != entry.ChangeType {
				continue
			}
			if rule.KeyPattern != "" && !globMatch(rule.KeyPattern, entry.Key) {
				continue
			}
			alerts = append(alerts, Alert{
				Rule:        rule.Name,
				Severity:    rule.Severity,
				Key:         entry.Key,
				ChangeType:  entry.ChangeType,
				TriggeredAt: now,
			})
		}
	}
	return alerts
}

// WriteAlerts writes alerts to w in the requested format ("text" or "json").
func WriteAlerts(w io.Writer, alerts []Alert, format string) error {
	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(alerts)
	default:
		for _, a := range alerts {
			fmt.Fprintf(w, "[%s] %s — key=%q change=%s at=%s\n",
				a.Severity, a.Rule, a.Key, a.ChangeType,
				a.TriggeredAt.Format(time.RFC3339))
		}
		return nil
	}
}
