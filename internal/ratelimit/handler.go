package ratelimit

import (
	"encoding/json"
	"net/http"
	"time"
)

type statusResponse struct {
	Container  string    `json:"container"`
	Allowed    bool      `json:"allowed"`
	LastSeen   time.Time `json:"last_seen,omitempty"`
	WindowSecs float64   `json:"window_secs"`
}

// Handler returns an HTTP handler that exposes rate-limit status and
// allows operators to reset the limiter for a specific container.
//
// GET  /ratelimit?container=<name>  — query current status
// POST /ratelimit/reset?container=<name> — reset the limiter
func Handler(rl *Limiter) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ratelimit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		container := r.URL.Query().Get("container")
		if container == "" {
			http.Error(w, "container query param required", http.StatusBadRequest)
			return
		}

		allowed := rl.Allow(container)
		// Undo the side-effect — we only wanted to inspect, so reset if it
		// was the first probe (allowed==true means a slot was consumed).
		if allowed {
			rl.Reset(container)
		}

		resp := statusResponse{
			Container:  container,
			Allowed:    allowed,
			WindowSecs: rl.window.Seconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	})

	mux.HandleFunc("/ratelimit/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		container := r.URL.Query().Get("container")
		if container == "" {
			http.Error(w, "container query param required", http.StatusBadRequest)
			return
		}
		rl.Reset(container)
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}
