// Package rollup aggregates drift results across multiple containers
// into a summarised report suitable for dashboards and alerting.
package rollup

import (
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Summary holds aggregated drift statistics for a single check cycle.
type Summary struct {
	Timestamp     time.Time
	Total         int
	Drifted       int
	Clean         int
	DriftRate     float64 // fraction 0.0–1.0
	ByContainer   map[string]ContainerSummary
}

// ContainerSummary holds per-container drift detail.
type ContainerSummary struct {
	Name    string
	Drifted bool
	Fields  []string // names of drifted fields
}

// Aggregator computes rollup summaries from drift results.
type Aggregator struct{}

// New returns a new Aggregator.
func New() *Aggregator {
	return &Aggregator{}
}

// Compute builds a Summary from a slice of drift results.
func (a *Aggregator) Compute(results []drift.Result) Summary {
	s := Summary{
		Timestamp:   time.Now().UTC(),
		Total:       len(results),
		ByContainer: make(map[string]ContainerSummary, len(results)),
	}

	for _, r := range results {
		fields := driftedFields(r)
		drifted := len(fields) > 0
		if drifted {
			s.Drifted++
		} else {
			s.Clean++
		}
		s.ByContainer[r.ContainerName] = ContainerSummary{
			Name:    r.ContainerName,
			Drifted: drifted,
			Fields:  fields,
		}
	}

	if s.Total > 0 {
		s.DriftRate = float64(s.Drifted) / float64(s.Total)
	}
	return s
}

// driftedFields returns the names of fields that have drifted in a result.
func driftedFields(r drift.Result) []string {
	var fields []string
	if r.ImageDrift {
		fields = append(fields, "image")
	}
	if len(r.EnvDrift) > 0 {
		fields = append(fields, "env")
	}
	if len(r.LabelDrift) > 0 {
		fields = append(fields, "labels")
	}
	return fields
}
