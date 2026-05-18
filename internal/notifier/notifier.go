package notifier

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/yourusername/driftwatch/internal/drift"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Notifier sends drift alerts to a configured output.
type Notifier struct {
	out       io.Writer
	threshold int // minimum number of drifted fields to trigger a notification
}

// New creates a Notifier writing to out.
// threshold sets the minimum drift count before a notification is emitted.
func New(out io.Writer, threshold int) *Notifier {
	if out == nil {
		out = os.Stderr
	}
	if threshold < 1 {
		threshold = 1
	}
	return &Notifier{out: out, threshold: threshold}
}

// Notify evaluates drift results and writes a notification line when the
// number of drifted containers meets or exceeds the configured threshold.
func (n *Notifier) Notify(results []drift.Result) error {
	drifted := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if r.HasDrift() {
			drifted = append(drifted, r)
		}
	}

	if len(drifted) < n.threshold {
		return nil
	}

	level := LevelWarn
	if len(drifted) >= n.threshold*2 {
		level = LevelError
	}

	names := make([]string, 0, len(drifted))
	for _, r := range drifted {
		names = append(names, r.ContainerName)
	}

	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s drift detected in %d container(s): %s\n",
		time.Now().UTC().Format(time.RFC3339),
		level,
		len(drifted),
		strings.Join(names, ", "),
	)
	return err
}
