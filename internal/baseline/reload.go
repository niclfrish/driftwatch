package baseline

import (
	"fmt"
	"log/slog"

	"github.com/yourorg/driftwatch/internal/manifest"
)

// Reloader watches a manifest file path and reloads the baseline Store
// whenever the manifest is re-read (e.g. on a scheduler tick).
type Reloader struct {
	path  string
	store *Store
	log   *slog.Logger
}

// NewReloader creates a Reloader that will update store from the manifest at path.
func NewReloader(path string, store *Store, log *slog.Logger) *Reloader {
	return &Reloader{path: path, store: store, log: log}
}

// Reload reads the manifest file and updates the store.
// It returns an error if the manifest cannot be loaded or is invalid.
func (r *Reloader) Reload() error {
	m, err := manifest.Load(r.path)
	if err != nil {
		return fmt.Errorf("baseline reloader: %w", err)
	}
	if err := r.store.Load(m); err != nil {
		return fmt.Errorf("baseline reloader: %w", err)
	}
	r.log.Info("baseline reloaded", "path", r.path, "containers", r.store.Len())
	return nil
}
