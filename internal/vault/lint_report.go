package vault

import (
	"fmt"
	"io"
	"strings"
)

// LintReport aggregates lint results across multiple secret paths.
type LintReport struct {
	Results []LintResult
}

// NewLintReport builds a LintReport from a map of path -> secrets using the given rules.
func NewLintReport(secrets map[string]map[string]string, rules []LintRule) *LintReport {
	report := &LintReport{}
	for path, kv := range secrets {
		result := LintSecrets(path, kv, rules)
		report.Results = append(report.Results, result)
	}
	return report
}

// TotalViolations returns the total number of violations across all results.
func (r *LintReport) TotalViolations() int {
	total := 0
	for _, res := range r.Results {
		total += len(res.Violations)
	}
	return total
}

// HasViolations returns true if any result has violations.
func (r *LintReport) HasViolations() bool {
	return r.TotalViolations() > 0
}

// Write outputs the lint report to the given writer.
func (r *LintReport) Write(w io.Writer) error {
	if len(r.Results) == 0 {
		_, err := fmt.Fprintln(w, "[lint] no secrets to lint")
		return err
	}
	var sb strings.Builder
	for _, res := range r.Results {
		sb.WriteString(res.Summary())
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("\n[lint] total violations: %d", r.TotalViolations()))
	_, err := fmt.Fprintln(w, sb.String())
	return err
}
