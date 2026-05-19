package audit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/driftwatch/internal/audit"
)

func TestStore_AppendAndQuery(t *testing.T) {
	s := audit.NewStore(10)
	s.Append(audit.Event{Kind: audit.EventDriftDetected, Container: "a"})
	s.Append(audit.Event{Kind: audit.EventScanCompleted, Container: "b"})

	all := s.Query(audit.QueryParams{})
	if len(all) != 2 {
		t.Fatalf("expected 2 events, got %d", len(all))
	}
}

func TestStore_QueryFilterContainer(t *testing.T) {
	s := audit.NewStore(10)
	s.Append(audit.Event{Kind: audit.EventDriftDetected, Container: "alpha"})
	s.Append(audit.Event{Kind: audit.EventDriftDetected, Container: "beta"})

	res := s.Query(audit.QueryParams{Container: "alpha"})
	if len(res) != 1 || res[0].Container != "alpha" {
		t.Errorf("expected 1 alpha event, got %v", res)
	}
}

func TestStore_CapEviction(t *testing.T) {
	s := audit.NewStore(3)
	for i := 0; i < 5; i++ {
		s.Append(audit.Event{Kind: audit.EventScanStarted, Message: string(rune('A' + i))})
	}
	all := s.Query(audit.QueryParams{})
	if len(all) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(all))
	}
}

func TestHandler_ReturnsJSON(t *testing.T) {
	s := audit.NewStore(10)
	s.Append(audit.Event{Kind: audit.EventDriftDetected, Container: "web", Message: "drift"})

	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	w := httptest.NewRecorder()
	audit.Handler(s)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var events []audit.Event
	if err := json.NewDecoder(w.Body).Decode(&events); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestTeeLogger_AppendsToStore(t *testing.T) {
	var buf bytes.Buffer
	s := audit.NewStore(10)
	l := audit.NewTeeLogger(&buf, s)

	_ = l.Log(audit.EventRuleSuppressed, "svc", "suppressed", nil)

	events := s.Query(audit.QueryParams{})
	if len(events) != 1 {
		t.Fatalf("expected 1 stored event, got %d", len(events))
	}
	if events[0].Kind != audit.EventRuleSuppressed {
		t.Errorf("kind = %q", events[0].Kind)
	}
	if buf.Len() == 0 {
		t.Error("expected output written to buffer")
	}
}
