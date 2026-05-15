package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes drift results to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a new Reporter writing to the given writer.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes a human-readable or structured summary of drift results.
func (r *Reporter) Report(results []drift.Result) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(results)
	default:
		return r.writeText(results)
	}
}

func (r *Reporter) writeText(results []drift.Result) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(r.out, "DriftWatch Report — %s\n", timestamp)
	fmt.Fprintf(r.out, "%s\n", strings.Repeat("-", 50))

	if len(results) == 0 {
		fmt.Fprintln(r.out, "No containers checked.")
		return nil
	}

	driftCount := 0
	for _, res := range results {
		if res.HasDrift {
			driftCount++
			fmt.Fprintf(r.out, "[DRIFT] %s\n", res.ContainerName)
			for _, d := range res.Differences {
				fmt.Fprintf(r.out, "  • %s\n", d)
			}
		} else {
			fmt.Fprintf(r.out, "[OK]    %s\n", res.ContainerName)
		}
	}

	fmt.Fprintf(r.out, "%s\n", strings.Repeat("-", 50))
	fmt.Fprintf(r.out, "Summary: %d/%d containers drifted\n", driftCount, len(results))
	return nil
}

func (r *Reporter) writeJSON(results []drift.Result) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(r.out, "{\"timestamp\":%q,\"results\":[\n", timestamp)
	for i, res := range results {
		diffs := "[]"
		if len(res.Differences) > 0 {
			quoted := make([]string, len(res.Differences))
			for j, d := range res.Differences {
				quoted[j] = fmt.Sprintf("%q", d)
			}
			diffs = "[" + strings.Join(quoted, ",") + "]"
		}
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out, "  {\"container\":%q,\"drift\":%v,\"differences\":%s}%s\n",
			res.ContainerName, res.HasDrift, diffs, comma)
	}
	fmt.Fprintln(r.out, "]}")
	return nil
}
