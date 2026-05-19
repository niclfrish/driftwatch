// Package baseline provides a thread-safe in-memory store of the desired
// container state as declared in a driftwatch manifest.
//
// The store is loaded from a parsed manifest.Manifest and queried by the
// drift detector to compare running container state against expectations.
//
// Usage:
//
//	store := baseline.New()
//	if err := store.Load(m); err != nil { ... }
//	entry, ok := store.Get("my-container")
package baseline
