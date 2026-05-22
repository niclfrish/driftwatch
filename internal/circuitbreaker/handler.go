package circuitbreaker

import (
	"encoding/json"
	"net/http"
)

type statusEntry struct {
	Container string `json:"container"`
	State     string `json:"state"`
	Failures  int    `json:"failures"`
}

// Handler returns an http.HandlerFunc that exposes circuit breaker states
// for all known keys as a JSON array.
func Handler(b *Breaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b.mu.Lock()
		entries := make([]statusEntry, 0, len(b.states))
		for key, state := range b.states {
			entries = append(entries, statusEntry{
				Container: key,
				State:     state.String(),
				Failures:  b.failures[key],
			})
		}
		b.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}
}

// Reset clears all state for key, returning the circuit to closed with zero
// failures. It is exposed via DELETE /circuitbreaker/{container}.
func (b *Breaker) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.failures, key)
	delete(b.states, key)
	delete(b.openedAt, key)
}
