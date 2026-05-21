package metrics_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	m := metrics.New()
	if m.RunsTotal != 0 || m.DriftTotal != 0 || m.ContainersChecked != 0 {
		t.Fatal("expected zero-value metrics")
	}
}

func TestRecordRun_Increments(t *testing.T) {
	m := metrics.New()
	m.RecordRun(5, 2)
	m.RecordRun(3, 0)

	if m.RunsTotal != 2 {
		t.Errorf("RunsTotal: want 2, got %d", m.RunsTotal)
	}
	if m.ContainersChecked != 8 {
		t.Errorf("ContainersChecked: want 8, got %d", m.ContainersChecked)
	}
	if m.DriftTotal != 2 {
		t.Errorf("DriftTotal: want 2, got %d", m.DriftTotal)
	}
	if m.LastRunDriftCount != 0 {
		t.Errorf("LastRunDriftCount: want 0 after second run, got %d", m.LastRunDriftCount)
	}
}

func TestHandler_ContentType(t *testing.T) {
	m := metrics.New()
	m.RecordRun(4, 1)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	m.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("unexpected Content-Type: %s", ct)
	}
}

func TestHandler_ContainsMetrics(t *testing.T) {
	m := metrics.New()
	m.RecordRun(10, 3)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	m.Handler()(rec, req)

	body, _ := io.ReadAll(rec.Body)
	s := string(body)

	cases := []string{
		"driftwatch_runs_total 1",
		"driftwatch_drift_total 3",
		"driftwatch_containers_checked_total 10",
		"driftwatch_last_run_drift_count 3",
	}
	for _, want := range cases {
		if !strings.Contains(s, want) {
			t.Errorf("body missing %q\nfull body:\n%s", want, s)
		}
	}
}

func TestHandler_ZeroRuns(t *testing.T) {
	m := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	m.Handler()(rec, req)

	body, _ := io.ReadAll(rec.Body)
	if !strings.Contains(string(body), "driftwatch_runs_total 0") {
		t.Error("expected zero run count in output")
	}
}
