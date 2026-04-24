package vault

import (
	"fmt"
	"io"
	"strings"
)

// PolicyReport summarises the result of enforcing a policy against a diff.
type PolicyReport struct {
	PolicyName string
	Violations []PolicyViolation
	Checked    int
}

// BuildPolicyReport runs EnforcePolicy and wraps the result in a PolicyReport.
func BuildPolicyReport(entries []DiffEntry, p Policy) PolicyReport {
	return PolicyReport{
		PolicyName: p.Name,
		Violations: EnforcePolicy(entries, p),
		Checked:    len(entries),
	}
}

// Passed returns true when no violations were found.
func (r PolicyReport) Passed() bool {
	return len(r.Violations) == 0
}

// WriteTo writes a human-readable summary of the report to w.
func (r PolicyReport) WriteTo(w io.Writer) (int64, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Policy: %s\n", r.PolicyName))
	sb.WriteString(fmt.Sprintf("Checked: %d entries\n", r.Checked))
	if r.Passed() {
		sb.WriteString("Result: PASS — no violations\n")
	} else {
		sb.WriteString(fmt.Sprintf("Result: FAIL — %d violation(s)\n", len(r.Violations)))
		for _, v := range r.Violations {
			sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n", v.Change, v.Key, v.Reason))
		}
	}
	n, err := fmt.Fprint(w, sb.String())
	return int64(n), err
}
