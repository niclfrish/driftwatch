package healthcheck_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/healthcheck"
)

func TestLoggingMiddleware_LogsRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	checker := healthcheck.New()
	handler := healthcheck.LoggingMiddleware(logger, checker.Handler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	line := buf.String()
	if !strings.Contains(line, "healthcheck") {
		t.Fatalf("expected log line to contain 'healthcheck', got: %q", line)
	}
	if !strings.Contains(line, "200") {
		t.Fatalf("expected log line to contain status code 200, got: %q", line)
	}
}

func TestLoggingMiddleware_Logs503(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	checker := healthcheck.New()
	checker.RecordRun(errSentinel("forced error"))
	handler := healthcheck.LoggingMiddleware(logger, checker.Handler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if !strings.Contains(buf.String(), "503") {
		t.Fatalf("expected 503 in log, got: %q", buf.String())
	}
}

type errSentinel string

func (e errSentinel) Error() string { return string(e) }
