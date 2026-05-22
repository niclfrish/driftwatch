package healthcheck_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/driftwatch/internal/healthcheck"
)

func TestNew_InitiallyHealthy(t *testing.T) {
	c := healthcheck.New()
	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRecordRun_SuccessKeepsOK(t *testing.T) {
	c := healthcheck.New()
	c.RecordRun(nil)
	var s healthcheck.Status
	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if !s.OK {
		t.Fatal("expected OK=true after successful run")
	}
	if s.RunCount != 1 {
		t.Fatalf("expected run_count=1, got %d", s.RunCount)
	}
}

func TestRecordRun_ErrorSetsUnhealthy(t *testing.T) {
	c := healthcheck.New()
	c.RecordRun(errors.New("docker timeout"))
	var s healthcheck.Status
	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if s.OK {
		t.Fatal("expected OK=false after error")
	}
	if s.LastError != "docker timeout" {
		t.Fatalf("unexpected last_error: %q", s.LastError)
	}
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}
}

func TestRecordRun_RecoveryAfterError(t *testing.T) {
	c := healthcheck.New()
	c.RecordRun(errors.New("boom"))
	c.RecordRun(nil)
	var s healthcheck.Status
	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if !s.OK {
		t.Fatal("expected recovery after successful run")
	}
	if s.LastError != "" {
		t.Fatalf("expected empty last_error, got %q", s.LastError)
	}
}

func TestHandler_ContentType(t *testing.T) {
	c := healthcheck.New()
	rec := httptest.NewRecorder()
	c.Handler()(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}
