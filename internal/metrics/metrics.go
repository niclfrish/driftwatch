// Package metrics exposes Prometheus-style counters and gauges for
// drift detection runs, making it easy to scrape operational data.
package metrics

import (
	"fmt"
	"net/http"
	"sync"
)

// Metrics holds counters and gauges collected during drift detection.
type Metrics struct {
	mu sync.RWMutex

	RunsTotal      int64
	DriftTotal     int64
	ContainersChecked int64
	LastRunDriftCount int
}

// New returns a zero-value Metrics instance.
func New() *Metrics {
	return &Metrics{}
}

// RecordRun increments the run counter and records how many containers
// were checked and how many had drift in the most recent run.
func (m *Metrics) RecordRun(containersChecked, driftCount int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RunsTotal++
	m.ContainersChecked += int64(containersChecked)
	m.DriftTotal += int64(driftCount)
	m.LastRunDriftCount = driftCount
}

// Handler returns an http.HandlerFunc that renders metrics in a simple
// Prometheus text exposition format.
func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "# HELP driftwatch_runs_total Total number of drift detection runs.\n")
		fmt.Fprintf(w, "# TYPE driftwatch_runs_total counter\n")
		fmt.Fprintf(w, "driftwatch_runs_total %d\n", m.RunsTotal)

		fmt.Fprintf(w, "# HELP driftwatch_drift_total Total containers found drifted across all runs.\n")
		fmt.Fprintf(w, "# TYPE driftwatch_drift_total counter\n")
		fmt.Fprintf(w, "driftwatch_drift_total %d\n", m.DriftTotal)

		fmt.Fprintf(w, "# HELP driftwatch_containers_checked_total Total container checks performed.\n")
		fmt.Fprintf(w, "# TYPE driftwatch_containers_checked_total counter\n")
		fmt.Fprintf(w, "driftwatch_containers_checked_total %d\n", m.ContainersChecked)

		fmt.Fprintf(w, "# HELP driftwatch_last_run_drift_count Number of drifted containers in the most recent run.\n")
		fmt.Fprintf(w, "# TYPE driftwatch_last_run_drift_count gauge\n")
		fmt.Fprintf(w, "driftwatch_last_run_drift_count %d\n", m.LastRunDriftCount)
	}
}
