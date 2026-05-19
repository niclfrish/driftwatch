package audit

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// QueryParams holds optional filters for listing audit events from a store.
type QueryParams struct {
	Container string
	Kind      EventKind
	Limit     int
}

// Store is an in-memory ring buffer of recent audit events.
type Store struct {
	cap    int
	events []Event
}

// NewStore creates a Store that retains at most capacity events.
func NewStore(capacity int) *Store {
	if capacity <= 0 {
		capacity = 200
	}
	return &Store{cap: capacity}
}

// Append adds an event to the store, evicting the oldest if at capacity.
func (s *Store) Append(ev Event) {
	if len(s.events) >= s.cap {
		s.events = s.events[1:]
	}
	s.events = append(s.events, ev)
}

// Query returns events matching the given parameters.
func (s *Store) Query(p QueryParams) []Event {
	out := make([]Event, 0, len(s.events))
	for _, ev := range s.events {
		if p.Container != "" && ev.Container != p.Container {
			continue
		}
		if p.Kind != "" && ev.Kind != p.Kind {
			continue
		}
		out = append(out, ev)
	}
	if p.Limit > 0 && len(out) > p.Limit {
		out = out[len(out)-p.Limit:]
	}
	return out
}

// Handler returns an http.HandlerFunc that serves recent audit events as JSON.
func Handler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		p := QueryParams{
			Container: q.Get("container"),
			Kind:      EventKind(q.Get("kind")),
		}
		if lStr := q.Get("limit"); lStr != "" {
			if n, err := strconv.Atoi(lStr); err == nil {
				p.Limit = n
			}
		}
		events := store.Query(p)
		if events == nil {
			events = []Event{}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	}
}

// TeeLogger wraps a Logger and also appends events to a Store.
type TeeLogger struct {
	*Logger
	store *Store
}

// NewTeeLogger creates a TeeLogger that writes to w and stores in s.
func NewTeeLogger(w interface{ Write([]byte) (int, error) }, s *Store) *TeeLogger {
	return &TeeLogger{Logger: New(w), store: s}
}

// Log emits an event to both the writer and the store.
func (t *TeeLogger) Log(kind EventKind, container, message string, meta map[string]string) error {
	ev := Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Container: container,
		Message:   message,
		Meta:      meta,
	}
	t.store.Append(ev)
	return t.Logger.Log(kind, container, message, meta)
}
