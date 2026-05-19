// Package audit provides structured audit logging for drift detection events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventDriftDetected  EventKind = "drift_detected"
	EventDriftResolved  EventKind = "drift_resolved"
	EventScanStarted    EventKind = "scan_started"
	EventScanCompleted  EventKind = "scan_completed"
	EventRuleSuppressed EventKind = "rule_suppressed"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Kind      EventKind         `json:"kind"`
	Container string            `json:"container,omitempty"`
	Message   string            `json:"message"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes audit events to an io.Writer as newline-delimited JSON.
type Logger struct {
	out io.Writer
}

// New creates a Logger that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Log emits an audit event.
func (l *Logger) Log(kind EventKind, container, message string, meta map[string]string) error {
	ev := Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Container: container,
		Message:   message,
		Meta:      meta,
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}
