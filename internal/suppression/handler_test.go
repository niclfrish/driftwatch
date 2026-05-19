package suppression_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/suppression"
)

func TestHandler_AddRule(t *testing.T) {
	s, _ := tempStore(t)
	h := suppression.Handler(s)

	body := `{"container":"web","field":"image","duration":"1h"}`
	req := httptest.NewRequest(http.MethodPost, "/suppressions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if !s.IsSuppressed("web", "image") {
		t.Error("expected rule to be active after add")
	}
}

func TestHandler_AddRule_MissingContainer(t *testing.T) {
	s, _ := tempStore(t)
	h := suppression.Handler(s)

	body := `{"duration":"1h"}`
	req := httptest.NewRequest(http.MethodPost, "/suppressions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_AddRule_BadDuration(t *testing.T) {
	s, _ := tempStore(t)
	h := suppression.Handler(s)

	body := `{"container":"web","duration":"notaduration"}`
	req := httptest.NewRequest(http.MethodPost, "/suppressions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_Prune(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Add(suppression.Rule{Container: "old", Until: time.Now().Add(-time.Minute)})
	h := suppression.Handler(s)

	req := httptest.NewRequest(http.MethodDelete, "/suppressions", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if s.IsSuppressed("old", "") {
		t.Error("expected expired rule to be pruned")
	}
}

func TestHandler_AddRule_ResponseBody(t *testing.T) {
	s, _ := tempStore(t)
	h := suppression.Handler(s)

	body := `{"container":"api","field":"env","duration":"30m"}`
	req := httptest.NewRequest(http.MethodPost, "/suppressions", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var rule suppression.Rule
	if err := json.NewDecoder(rec.Body).Decode(&rule); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if rule.Container != "api" {
		t.Errorf("expected container 'api', got %q", rule.Container)
	}
	if rule.Until.Before(time.Now()) {
		t.Error("expected Until to be in the future")
	}
}
