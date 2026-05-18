package watcher_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/watcher"
)

// TestNew_ReturnsWatcher ensures New constructs a non-nil Watcher.
func TestNew_ReturnsWatcher(t *testing.T) {
	w := watcher.New("manifest.yaml", nil, nil, time.Second)
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

// TestRun_CancelImmediately verifies that Run exits cleanly when the context
// is cancelled before the first tick fires.
func TestRun_CancelImmediately(t *testing.T) {
	// Write a minimal valid manifest so the initial check can load it.
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "manifest.yaml")
	content := `containers:
  - name: web
    image: nginx:latest
`
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	// Use a very long interval so the ticker never fires during the test.
	w := watcher.New(manifestPath, nil, nil, 10*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan error, 1)
	go func() {
		done <- w.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}
