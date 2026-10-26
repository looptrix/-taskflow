package worker

import (
	"context"
	"testing"
	"time"
)

func TestSubmitExecutesJob(t *testing.T) {
	pool := NewPool(context.Background(), 2, 10)
	defer pool.Close()

	done := make(chan struct{})
	if !pool.Submit(func(ctx context.Context) error {
		close(done)
		return nil
	}) {
		t.Fatal("expected job to be accepted")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("job was not executed")
	}
}

func TestSubmitRejectsNilJob(t *testing.T) {
	pool := NewPool(context.Background(), 1, 1)
	defer pool.Close()

	if pool.Submit(nil) {
		t.Fatal("nil job must be rejected")
	}
}

func TestPoolNormalizesWorkerCount(t *testing.T) {
	pool := NewPool(context.Background(), 0, 0)
	defer pool.Close()

	done := make(chan struct{})
	if !pool.Submit(func(ctx context.Context) error {
		close(done)
		return nil
	}) {
		t.Fatal("pool should accept a job with zero configuration")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("job was not executed")
	}
}
