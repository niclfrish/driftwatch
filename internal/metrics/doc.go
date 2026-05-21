// Package metrics provides a lightweight in-process metrics store and an
// HTTP handler that exposes counters and gauges in Prometheus text format.
//
// Usage:
//
//	m := metrics.New()
//
//	// after each drift detection run:
//	m.RecordRun(len(containers), len(driftedContainers))
//
//	// mount the scrape endpoint:
//	http.Handle("/metrics", m.Handler())
//
// No external dependencies are required; the exposition format is a
// minimal subset of the Prometheus text protocol (version 0.0.4).
package metrics
