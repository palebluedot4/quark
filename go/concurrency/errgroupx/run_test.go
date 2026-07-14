package errgroupx_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/errgroupx"
)

func TestRunAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		n         int
		taskCount int
	}{
		{
			name:      "single worker",
			n:         1,
			taskCount: 10,
		},
		{
			name:      "bounded below task count",
			n:         4,
			taskCount: 20,
		},
		{
			name:      "limit exceeds task count",
			n:         10,
			taskCount: 3,
		},
		{
			name:      "empty tasks",
			n:         4,
			taskCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var ran, current, peak atomic.Int64
			tasks := make([]func(context.Context) error, tt.taskCount)
			for i := range tasks {
				tasks[i] = func(context.Context) error {
					ran.Add(1)
					c := current.Add(1)
					for {
						p := peak.Load()
						if c <= p || peak.CompareAndSwap(p, c) {
							break
						}
					}
					time.Sleep(time.Millisecond)
					current.Add(-1)
					return nil
				}
			}
			if err := errgroupx.RunAll(context.Background(), tt.n, tasks); err != nil {
				t.Errorf("RunAll() = %v, want nil", err)
			}
			if got := ran.Load(); got != int64(tt.taskCount) {
				t.Errorf("RunAll() executed tasks = %v, want %v", got, tt.taskCount)
			}
			if got := peak.Load(); got > int64(tt.n) {
				t.Errorf("RunAll() peak concurrency = %v, want <= %v", got, tt.n)
			}
		})
	}
}

func TestRunAllReturnsFirstErrorAndCancels(t *testing.T) {
	t.Parallel()
	want := errors.New("boom")
	var cancelled atomic.Int64
	tasks := []func(context.Context) error{
		func(context.Context) error {
			return want
		},
	}
	const cooperating = 5
	for range cooperating {
		tasks = append(tasks, func(ctx context.Context) error {
			<-ctx.Done()
			cancelled.Add(1)
			return ctx.Err()
		})
	}
	if err := errgroupx.RunAll(context.Background(), len(tasks), tasks); !errors.Is(err, want) {
		t.Errorf("RunAll() = %v, want %v", err, want)
	}
	if got := cancelled.Load(); got != cooperating {
		t.Errorf("RunAll() cancelled tasks = %v, want %v", got, cooperating)
	}
}

func TestRunAllRespectsParentCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tasks := make([]func(context.Context) error, 5)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error {
			return ctx.Err()
		}
	}
	if err := errgroupx.RunAll(ctx, 2, tasks); !errors.Is(err, context.Canceled) {
		t.Errorf("RunAll() = %v, want %v", err, context.Canceled)
	}
}

func TestRunAllPanicsOnNonPositiveLimit(t *testing.T) {
	t.Parallel()
	for _, n := range []int{0, -1} {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			t.Parallel()
			defer func() {
				if recover() == nil {
					t.Errorf("RunAll(%d, nil) did not panic, want panic", n)
				}
			}()
			_ = errgroupx.RunAll(context.Background(), n, nil)
		})
	}
}

func BenchmarkRunAll(b *testing.B) {
	ctx := context.Background()
	tasks := make([]func(context.Context) error, 1000)
	for i := range tasks {
		tasks[i] = func(context.Context) error { return nil }
	}
	for b.Loop() {
		_ = errgroupx.RunAll(ctx, 10, tasks)
	}
}
