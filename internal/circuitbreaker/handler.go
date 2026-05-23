package circuitbreaker

import (
	"encoding/json"
	"net/http"
)

type circuitStatus struct {
	State    string `json:"state"`
	Failures int    `json:"failures"`
}

type statusResponse struct {
	Circuits map[string]circuitStatus `json:"circuits"`
}

// Handler returns an http.Handler that exposes the current state of all
// circuit breakers tracked by cb. An optional ?container= query parameter
// filters the response to a single circuit.
func Handler(cb *CircuitBreaker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("container")

		cb.mu.RLock()
		defer cb.mu.RUnlock()

		circuits := make(map[string]circuitStatus)
		for key, entry := range cb.entries {
			if filter != "" && key != filter {
				continue
			}
			state := "closed"
			switch entry.state {
			case stateOpen:
				state = "open"
			case stateHalfOpen:
				state = "half-open"
			}
			circuits[key] = circuitStatus{
				State:    state,
				Failures: entry.failures,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(statusResponse{Circuits: circuits})
	})
}
