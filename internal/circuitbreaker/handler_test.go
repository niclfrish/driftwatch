package circuitbreaker_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/driftwatch/internal/circuitbreaker"
)

func TestHandler_ContentType(t *testing.T) {
	cb := circuitbreaker.New(3, 30)
	h := circuitbreaker.Handler(cb)

	req := httptest.NewRequest(http.MethodGet, "/circuitbreaker", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected application/json, got %q", got)
	}
}

func TestHandler_ReturnsOK(t *testing.T) {
	cb := circuitbreaker.New(3, 30)
	h := circuitbreaker.Handler(cb)

	req := httptest.NewRequest(http.MethodGet, "/circuitbreaker", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_EmptyState(t *testing.T) {
	cb := circuitbreaker.New(3, 30)
	h := circuitbreaker.Handler(cb)

	req := httptest.NewRequest(http.MethodGet, "/circuitbreaker", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var result map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := result["circuits"]; !ok {
		t.Error("expected 'circuits' key in response")
	}
}

func TestHandler_ShowsOpenCircuit(t *testing.T) {
	cb := circuitbreaker.New(2, 30)
	for i := 0; i < 3; i++ {
		cb.RecordFailure("web")
	}
	h := circuitbreaker.Handler(cb)

	req := httptest.NewRequest(http.MethodGet, "/circuitbreaker", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var result struct {
		Circuits map[string]struct {
			State string `json:"state"`
		} `json:"circuits"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	entry, ok := result.Circuits["web"]
	if !ok {
		t.Fatal("expected 'web' circuit in response")
	}
	if entry.State != "open" {
		t.Errorf("expected state=open, got %q", entry.State)
	}
}

func TestHandler_FilterByContainer(t *testing.T) {
	cb := circuitbreaker.New(2, 30)
	cb.RecordFailure("alpha")
	cb.RecordFailure("beta")
	h := circuitbreaker.Handler(cb)

	req := httptest.NewRequest(http.MethodGet, "/circuitbreaker?container=alpha", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var result struct {
		Circuits map[string]interface{} `json:"circuits"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if _, ok := result.Circuits["beta"]; ok {
		t.Error("expected 'beta' to be filtered out")
	}
	if _, ok := result.Circuits["alpha"]; !ok {
		t.Error("expected 'alpha' to be present")
	}
}
