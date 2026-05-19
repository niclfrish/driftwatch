package suppression

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// addRequest is the JSON body for POST /suppressions.
type addRequest struct {
	Container string `json:"container"`
	Field     string `json:"field"`
	Duration  string `json:"duration"` // e.g. "2h", "30m"
}

// Handler returns an http.Handler that exposes suppression management
// endpoints under the given mux.
//
//	POST /suppressions  — add a rule
//	DELETE /suppressions — prune expired rules
func Handler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/suppressions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleAdd(s, w, r)
		case http.MethodDelete:
			handlePrune(s, w)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}

func handleAdd(s *Store, w http.ResponseWriter, r *http.Request) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}
	if req.Container == "" {
		http.Error(w, "container is required", http.StatusBadRequest)
		return
	}
	d, err := time.ParseDuration(req.Duration)
	if err != nil || d <= 0 {
		http.Error(w, "duration must be a positive Go duration string", http.StatusBadRequest)
		return
	}
	rule := Rule{
		Container: req.Container,
		Field:     req.Field,
		Until:     time.Now().Add(d),
	}
	if err := s.Add(rule); err != nil {
		http.Error(w, fmt.Sprintf("store error: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rule)
}

func handlePrune(s *Store, w http.ResponseWriter) {
	if err := s.Prune(); err != nil {
		http.Error(w, fmt.Sprintf("prune error: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
