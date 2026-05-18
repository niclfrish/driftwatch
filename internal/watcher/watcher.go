package watcher

import (
	"context"
	"log"
	"time"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/manifest"
	"github.com/user/driftwatch/internal/reporter"
)

// Watcher periodically checks for config drift between running containers
// and their source manifests, reporting any findings.
type Watcher struct {
	manifestPath string
	detector     *drift.Detector
	reporter     *reporter.Reporter
	interval     time.Duration
}

// New creates a new Watcher with the given manifest path, detector, reporter,
// and polling interval.
func New(manifestPath string, d *drift.Detector, r *reporter.Reporter, interval time.Duration) *Watcher {
	return &Watcher{
		manifestPath: manifestPath,
		detector:     d,
		reporter:     r,
		interval:     interval,
	}
}

// Run starts the watch loop, checking for drift at every interval tick.
// It blocks until the provided context is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	log.Printf("watcher: starting, interval=%s manifest=%s", w.interval, w.manifestPath)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run once immediately before waiting for the first tick.
	if err := w.check(ctx); err != nil {
		log.Printf("watcher: check error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("watcher: stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := w.check(ctx); err != nil {
				log.Printf("watcher: check error: %v", err)
			}
		}
	}
}

// check loads the manifest, inspects running containers, and reports drift.
func (w *Watcher) check(ctx context.Context) error {
	m, err := manifest.Load(w.manifestPath)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	results, err := w.detector.Detect(ctx, m)
	if err != nil {
		return fmt.Errorf("detect drift: %w", err)
	}

	return w.reporter.Report(results)
}
