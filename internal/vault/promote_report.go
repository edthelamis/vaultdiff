package vault

import (
	"fmt"
	"io"
	"time"
)

// PromoteReport captures metadata and results of a promotion operation.
type PromoteReport struct {
	Timestamp   time.Time
	Source      string
	Destination string
	Options     PromoteOptions
	Result      *PromoteResult
}

// NewPromoteReport creates a PromoteReport stamped with the current time.
func NewPromoteReport(src, dst string, opts PromoteOptions, result *PromoteResult) *PromoteReport {
	return &PromoteReport{
		Timestamp:   time.Now().UTC(),
		Source:      src,
		Destination: dst,
		Options:     opts,
		Result:      result,
	}
}

// Write renders the promotion report to w.
func (r *PromoteReport) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "Promote Report\n")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "  Time:        %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  Source:      %s\n", r.Source)
	fmt.Fprintf(w, "  Destination: %s\n", r.Destination)
	fmt.Fprintf(w, "  DryRun:      %v\n", r.Options.DryRun)
	fmt.Fprintf(w, "  Overwrite:   %v\n", r.Options.Overwrite)
	fmt.Fprintf(w, "  Summary:     %s\n", r.Result.Summary())

	if len(r.Result.Promoted) > 0 {
		fmt.Fprintf(w, "  Promoted:\n")
		for _, k := range r.Result.Promoted {
			fmt.Fprintf(w, "    + %s\n", k)
		}
	}
	if len(r.Result.Skipped) > 0 {
		fmt.Fprintf(w, "  Skipped:\n")
		for _, k := range r.Result.Skipped {
			fmt.Fprintf(w, "    ~ %s\n", k)
		}
	}
	if len(r.Result.Errors) > 0 {
		fmt.Fprintf(w, "  Errors:\n")
		for _, e := range r.Result.Errors {
			fmt.Fprintf(w, "    ! %s\n", e)
		}
	}
	return nil
}
