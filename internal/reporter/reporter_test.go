package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func makeResults(hasDrift bool) []drift.Result {
	if hasDrift {
		return []drift.Result{
			{
				ContainerName: "web",
				HasDrift:      true,
				Differences:   []string{"env APP_ENV: want=production got=staging"},
			},
			{
				ContainerName: "db",
				HasDrift:      false,
				Differences:   nil,
			},
		}
	}
	return []drift.Result{
		{ContainerName: "web", HasDrift: false},
		{ContainerName: "db", HasDrift: false},
	}
}

func TestReport_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	results := makeResults(false)
	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]    web") {
		t.Errorf("expected OK for web, got:\n%s", out)
	}
	if !strings.Contains(out, "0/2 containers drifted") {
		t.Errorf("expected summary 0/2, got:\n%s", out)
	}
}

func TestReport_TextWithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	results := makeResults(true)
	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DRIFT] web") {
		t.Errorf("expected DRIFT for web, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected difference detail, got:\n%s", out)
	}
	if !strings.Contains(out, "1/2 containers drifted") {
		t.Errorf("expected summary 1/2, got:\n%s", out)
	}
}

func TestReport_TextEmpty(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	if err := r.Report([]drift.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No containers checked") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatJSON)
	results := makeResults(true)
	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drift":true`) {
		t.Errorf("expected drift:true in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, `"container":"web"`) {
		t.Errorf("expected container name in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, "timestamp") {
		t.Errorf("expected timestamp field in JSON, got:\n%s", out)
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	r := New(nil, FormatText)
	if r.out == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
