package scheduler_test

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/scheduler"
)

func silentLogger() *log.Logger {
	return log.New(log.Writer(), "", 0)
}

func TestNew_ReturnsScheduler(t *testing.T) {
	job := func(ctx context.Context) error { return nil }
	s := scheduler.New(time.Second, job, nil)
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRun_ExecutesJobImmediately(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	s := scheduler.New(10*time.Second, job, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = s.Run(ctx)

	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected job to run at least once immediately")
	}
}

func TestRun_ExecutesJobOnTick(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	s := scheduler.New(30*time.Millisecond, job, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	_ = s.Run(ctx)

	if atomic.LoadInt32(&count) < 2 {
		t.Errorf("expected at least 2 executions, got %d", atomic.LoadInt32(&count))
	}
}

func TestRun_ContinuesOnJobError(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errors.New("job failed")
	}

	s := scheduler.New(30*time.Millisecond, job, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := s.Run(ctx)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
	if atomic.LoadInt32(&count) < 2 {
		t.Errorf("expected multiple runs despite errors, got %d", atomic.LoadInt32(&count))
	}
}

func TestRun_CancelReturnsCtxErr(t *testing.T) {
	job := func(ctx context.Context) error { return nil }
	s := scheduler.New(time.Second, job, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected Canceled, got %v", err)
	}
}
