package ratelimit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestLimiter() *Limiter {
	return New(30 * time.Second)
}

func TestHandler_GetStatus_MissingContainer(t *testing.T) {
	h := Handler(newTestLimiter())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_GetStatus_ReturnsJSON(t *testing.T) {
	h := Handler(newTestLimiter())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ratelimit?container=web", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp statusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Container != "web" {
		t.Errorf("expected container=web, got %q", resp.Container)
	}
	if resp.WindowSecs != 30 {
		t.Errorf("expected window_secs=30, got %f", resp.WindowSecs)
	}
}

func TestHandler_GetStatus_ContentType(t *testing.T) {
	h := Handler(newTestLimiter())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ratelimit?container=api", nil)
	h.ServeHTTP(rec, req)
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestHandler_Reset_NoContent(t *testing.T) {
	h := Handler(newTestLimiter())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ratelimit/reset?container=web", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestHandler_Reset_MissingContainer(t *testing.T) {
	h := Handler(newTestLimiter())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ratelimit/reset", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := Handler(newTestLimiter())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/ratelimit?container=web", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
