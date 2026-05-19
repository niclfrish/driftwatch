package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/audit"
)

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.EventDriftDetected, "app", "env mismatch", map[string]string{"field": "ENV_VAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ev audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &ev); err != nil {
		t.Fatalf("failed to parse output as JSON: %v", err)
	}
	if ev.Kind != audit.EventDriftDetected {
		t.Errorf("kind = %q, want %q", ev.Kind, audit.EventDriftDetected)
	}
	if ev.Container != "app" {
		t.Errorf("container = %q, want %q", ev.Container, "app")
	}
	if ev.Message != "env mismatch" {
		t.Errorf("message = %q, want %q", ev.Message, "env mismatch")
	}
	if ev.Meta["field"] != "ENV_VAR" {
		t.Errorf("meta field = %q, want %q", ev.Meta["field"], "ENV_VAR")
	}
	if ev.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestLog_NilMeta(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Log(audit.EventScanStarted, "", "scan begin", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), string(audit.EventScanStarted)) {
		t.Error("output should contain event kind")
	}
}

func TestLog_DefaultsToStdout(t *testing.T) {
	// Just ensure New(nil) doesn't panic.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	kinds := []audit.EventKind{
		audit.EventScanStarted,
		audit.EventDriftDetected,
		audit.EventScanCompleted,
	}
	for _, k := range kinds {
		if err := l.Log(k, "c1", "msg", nil); err != nil {
			t.Fatalf("Log(%q): %v", k, err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	for i, line := range lines {
		var ev audit.Event
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			t.Errorf("line %d: %v", i, err)
		}
		if ev.Kind != kinds[i] {
			t.Errorf("line %d kind = %q, want %q", i, ev.Kind, kinds[i])
		}
	}
}
