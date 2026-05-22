package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/your-org/driftwatch/internal/circuitbreaker"
)

func TestAllow_InitiallyClosed(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if !b.Allow("c1") {
		t.Fatal("expected Allow=true for fresh key")
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		b.RecordFailure("c1")
	}
	if b.Allow("c1") {
		t.Fatal("expected Allow=false after threshold failures")
	}
	if b.State("c1") != circuitbreaker.StateOpen {
		t.Fatalf("expected StateOpen, got %s", b.State("c1"))
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b := circuitbreaker.New(2, time.Minute)
	b.RecordFailure("c1")
	b.RecordFailure("c1")
	b.RecordSuccess("c1")
	if !b.Allow("c1") {
		t.Fatal("expected Allow=true after success reset")
	}
	if b.State("c1") != circuitbreaker.StateClosed {
		t.Fatalf("expected StateClosed, got %s", b.State("c1"))
	}
}

func TestHalfOpen_AfterResetDuration(t *testing.T) {
	b := circuitbreaker.New(1, 50*time.Millisecond)
	b.RecordFailure("c1")
	if b.Allow("c1") {
		t.Fatal("expected Allow=false immediately after open")
	}
	time.Sleep(60 * time.Millisecond)
	if !b.Allow("c1") {
		t.Fatal("expected Allow=true after reset duration (half-open probe)")
	}
	if b.State("c1") != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %s", b.State("c1"))
	}
}

func TestIndependentKeys(t *testing.T) {
	b := circuitbreaker.New(2, time.Minute)
	b.RecordFailure("c1")
	b.RecordFailure("c1")
	if !b.Allow("c2") {
		t.Fatal("c2 should be unaffected by c1 failures")
	}
}

func TestStateString(t *testing.T) {
	cases := []struct {
		s    circuitbreaker.State
		want string
	}{
		{circuitbreaker.StateClosed, "closed"},
		{circuitbreaker.StateOpen, "open"},
		{circuitbreaker.StateHalfOpen, "half-open"},
	}
	for _, tc := range cases {
		if tc.s.String() != tc.want {
			t.Errorf("State.String() = %q, want %q", tc.s.String(), tc.want)
		}
	}
}
