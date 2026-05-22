// Package history provides persistent storage for drift detection results.
//
// Each call to [Store.Record] writes a timestamped JSON file to a configured
// directory. Entries can be retrieved in chronological order via [Store.List]
// or the most recent entry via [Store.Latest], enabling trend analysis and
// change detection between successive driftwatch runs.
//
// # Directory Layout
//
// Files are stored under the configured root directory using the naming
// convention "<RFC3339-timestamp>.json", for example:
//
//	/var/lib/driftwatch/history/
//		2024-01-15T10:30:00Z.json
//		2024-01-15T11:00:00Z.json
//		2024-01-15T11:30:00Z.json
//
// # Concurrency
//
// Store is safe for concurrent use. Each [Store.Record] call generates a
// unique filename based on the current time, so simultaneous writes from
// multiple goroutines will not collide under normal clock conditions.
package history
