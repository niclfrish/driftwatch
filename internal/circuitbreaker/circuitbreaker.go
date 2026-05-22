// Package circuitbreaker provides a simple per-container circuit breaker
// that stops drift checks when a container endpoint is repeatedly unreachable.
package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// State represents the current state of a circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing; requests blocked
	StateHalfOpen              // probe allowed
)

// Breaker is a per-key circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	failures    map[string]int
	states      map[string]State
	openedAt    map[string]time.Time
	threshold   int
	resetAfter  time.Duration
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts a half-open probe after resetAfter duration.
func New(threshold int, resetAfter time.Duration) *Breaker {
	return &Breaker{
		failures:   make(map[string]int),
		states:     make(map[string]State),
		openedAt:   make(map[string]time.Time),
		threshold:  threshold,
		resetAfter: resetAfter,
	}
}

// Allow returns true when the circuit is closed or half-open for the given key.
func (b *Breaker) Allow(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.states[key] {
	case StateOpen:
		if time.Since(b.openedAt[key]) >= b.resetAfter {
			b.states[key] = StateHalfOpen
			return true
		}
		return false
	default:
		return true
	}
}

// RecordSuccess resets the failure count and closes the circuit for key.
func (b *Breaker) RecordSuccess(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[key] = 0
	b.states[key] = StateClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (b *Breaker) RecordFailure(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[key]++
	if b.failures[key] >= b.threshold {
		b.states[key] = StateOpen
		b.openedAt[key] = time.Now()
	}
}

// State returns the current circuit state for key.
func (b *Breaker) State(key string) State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.states[key]
}

// String implements fmt.Stringer for State.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown(%d)", int(s))
	}
}
