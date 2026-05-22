package ratelimit_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/ratelimit"
)

const window = 5 * time.Minute

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l := ratelimit.New(window)
	if !l.Allow("web") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinWindowDenied(t *testing.T) {
	l := ratelimit.New(window)
	now := time.Now()
	l.AllowAt("web", now)
	if l.AllowAt("web", now.Add(time.Minute)) {
		t.Fatal("expected second call within window to be denied")
	}
}

func TestAllow_CallAfterWindowAllowed(t *testing.T) {
	l := ratelimit.New(window)
	now := time.Now()
	l.AllowAt("web", now)
	if !l.AllowAt("web", now.Add(window)) {
		t.Fatal("expected call exactly at window boundary to be allowed")
	}
}

func TestAllow_IndependentContainers(t *testing.T) {
	l := ratelimit.New(window)
	now := time.Now()
	l.AllowAt("web", now)
	if !l.AllowAt("api", now.Add(time.Second)) {
		t.Fatal("expected different container to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	l := ratelimit.New(window)
	now := time.Now()
	l.AllowAt("web", now)
	l.Reset("web")
	if !l.AllowAt("web", now.Add(time.Second)) {
		t.Fatal("expected reset container to be allowed immediately")
	}
}

func TestPrune_RemovesStaleEntries(t *testing.T) {
	l := ratelimit.New(window)
	now := time.Now()
	l.AllowAt("web", now)
	// Prune at a time beyond the window; subsequent call should be allowed
	l.Prune(now.Add(window + time.Second))
	if !l.AllowAt("web", now.Add(window+time.Second)) {
		t.Fatal("expected pruned container to be allowed again")
	}
}

func TestNew_DefaultsNonPositiveRate(t *testing.T) {
	l := ratelimit.New(0)
	now := time.Now()
	l.AllowAt("x", now)
	// Default rate is 1 minute; 30s later should be denied
	if l.AllowAt("x", now.Add(30*time.Second)) {
		t.Fatal("expected call within default 1-minute window to be denied")
	}
}
