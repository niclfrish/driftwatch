// Package healthcheck exposes a simple HTTP health endpoint that reports
// the overall readiness of the driftwatch daemon.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the current health of the daemon.
type Status struct {
	OK        bool      `json:"ok"`
	LastRun   time.Time `json:"last_run,omitempty"`
	RunCount  int64     `json:"run_count"`
	LastError string    `json:"last_error,omitempty"`
}

// Checker holds daemon health state and serves it over HTTP.
type Checker struct {
	mu     sync.RWMutex
	status Status
}

// New returns an initialised Checker.
func New() *Checker {
	return &Checker{
		status: Status{OK: true},
	}
}

// RecordRun updates the health state after each scheduler tick.
func (c *Checker) RecordRun(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.LastRun = time.Now().UTC()
	c.status.RunCount++
	if err != nil {
		c.status.OK = false
		c.status.LastError = err.Error()
	} else {
		c.status.OK = true
		c.status.LastError = ""
	}
}

// Handler returns an http.HandlerFunc that writes the current Status as JSON.
func (c *Checker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.RLock()
		s := c.status
		c.mu.RUnlock()

		code := http.StatusOK
		if !s.OK {
			code = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(s)
	}
}
