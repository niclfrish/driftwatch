package notifier_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/driftwatch/internal/drift"
	"github.com/yourusername/driftwatch/internal/notifier"
)

func makeResult(name string, diffs []drift.FieldDiff) drift.Result {
	return drift.Result{
		ContainerName: name,
		Diffs:         diffs,
	}
}

func TestNotify_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf, 1)

	err := n.Notify([]drift.Result{makeResult("app", nil)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %q", buf.String())
	}
}

func TestNotify_BelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf, 3)

	diffs := []drift.FieldDiff{{Field: "image", Expected: "a", Actual: "b"}}
	results := []drift.Result{
		makeResult("svc1", diffs),
		makeResult("svc2", diffs),
	}

	if err := n.Notify(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output below threshold, got: %q", buf.String())
	}
}

func TestNotify_WarnLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf, 1)

	diffs := []drift.FieldDiff{{Field: "image", Expected: "v1", Actual: "v2"}}
	if err := n.Notify([]drift.Result{makeResult("worker", diffs)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %q", out)
	}
	if !strings.Contains(out, "worker") {
		t.Errorf("expected container name in output, got: %q", out)
	}
}

func TestNotify_ErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf, 1)

	diffs := []drift.FieldDiff{{Field: "env.PORT", Expected: "8080", Actual: "9090"}}
	results := []drift.Result{
		makeResult("api", diffs),
		makeResult("proxy", diffs),
		makeResult("cache", diffs),
	}

	if err := n.Notify(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR level for high drift count, got: %q", out)
	}
}

func TestNew_DefaultThreshold(t *testing.T) {
	n := notifier.New(nil, 0) // threshold < 1 should default to 1
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
