// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently drift notifications are emitted per container.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-container notification rate using a token-bucket approach.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]bucket
	rate     time.Duration // minimum duration between allowed events
}

type bucket struct {
	lastAllowed time.Time
}

// New creates a Limiter that allows at most one notification per container
// within the given rate window (e.g. 5*time.Minute).
func New(rate time.Duration) *Limiter {
	if rate <= 0 {
		rate = time.Minute
	}
	return &Limiter{
		buckets: make(map[string]bucket),
		rate:    rate,
	}
}

// Allow reports whether a notification for the given container should be
// allowed at the current time. If allowed, the internal timestamp is updated.
func (l *Limiter) Allow(containerName string) bool {
	return l.AllowAt(containerName, time.Now())
}

// AllowAt is like Allow but accepts an explicit time, useful for testing.
func (l *Limiter) AllowAt(containerName string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.buckets[containerName]
	if !exists || now.Sub(b.lastAllowed) >= l.rate {
		l.buckets[containerName] = bucket{lastAllowed: now}
		return true
	}
	return false
}

// Reset clears the rate-limit state for a specific container.
func (l *Limiter) Reset(containerName string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, containerName)
}

// Prune removes entries whose last-allowed time is older than the rate window,
// keeping memory usage bounded over long runtimes.
func (l *Limiter) Prune(now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for name, b := range l.buckets {
		if now.Sub(b.lastAllowed) >= l.rate {
			delete(l.buckets, name)
		}
	}
}
