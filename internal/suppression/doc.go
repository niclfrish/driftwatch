// Package suppression manages drift alert suppression rules.
//
// Rules can target a specific container and optionally a specific field
// (e.g. "image", "env"). A rule with an empty Field suppresses all drift
// for the named container until the rule expires.
//
// Rules are persisted to a JSON file so they survive daemon restarts.
// Call Prune periodically (e.g. via the scheduler) to evict expired entries.
package suppression
