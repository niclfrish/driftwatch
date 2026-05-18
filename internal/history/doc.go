// Package history provides persistent storage for drift detection results.
//
// Each call to [Store.Record] writes a timestamped JSON file to a configured
// directory. Entries can be retrieved in chronological order via [Store.List]
// or the most recent entry via [Store.Latest], enabling trend analysis and
// change detection between successive driftwatch runs.
package history
