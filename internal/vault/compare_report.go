package vault

import (
	"fmt"
	"io"
	"strings"
)

// CompareReport holds a formatted comparison between two environments.
type CompareReport struct {
	EnvA   string
	EnvB   string
	Result CompareResult
}

// NewCompareReport creates a CompareReport for two named environments.
func NewCompareReport(envA, envB string, result CompareResult) *CompareReport {
	return &CompareReport{
		EnvA:   envA,
		EnvB:   envB,
		Result: result,
	}
}

// Write renders the report to the given writer.
func (r *CompareReport) Write(w io.Writer) error {
	fmt.Fprintf(w, "=== Secret Comparison: %s vs %s ===\n", r.EnvA, r.EnvB)

	if len(r.Result.Matching) > 0 {
		fmt.Fprintf(w, "\n[MATCHING] (%d)\n", len(r.Result.Matching))
		for _, k := range r.Result.Matching {
			fmt.Fprintf(w, "  = %s\n", k)
		}
	}

	if len(r.Result.Differing) > 0 {
		fmt.Fprintf(w, "\n[DIFFERING] (%d)\n", len(r.Result.Differing))
		for _, k := range r.Result.Differing {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	if len(r.Result.OnlyInA) > 0 {
		fmt.Fprintf(w, "\n[ONLY IN %s] (%d)\n", strings.ToUpper(r.EnvA), len(r.Result.OnlyInA))
		for _, k := range r.Result.OnlyInA {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.Result.OnlyInB) > 0 {
		fmt.Fprintf(w, "\n[ONLY IN %s] (%d)\n", strings.ToUpper(r.EnvB), len(r.Result.OnlyInB))
		for _, k := range r.Result.OnlyInB {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	fmt.Fprintln(w)
	return nil
}
