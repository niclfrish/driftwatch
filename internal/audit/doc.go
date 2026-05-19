// Package audit provides structured, append-only audit logging for
// driftwatch events. Each event is written as a newline-delimited JSON
// record so that logs can be consumed by standard tooling such as jq,
// Loki, or Elasticsearch.
//
// Usage:
//
//	logger := audit.New(os.Stderr)
//	logger.Log(audit.EventDriftDetected, "my-container", "env var changed", nil)
package audit
