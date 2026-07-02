package semaphorex_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/semaphorex"
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
			tasks := make([]func(), tt.taskCount)
			for i := range tasks {
				tasks[i] = func() {
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
				}
			}
			semaphorex.RunAll(tt.n, tasks)
			if got := ran.Load(); got != int64(tt.taskCount) {
				t.Errorf("RunAll() executed tasks = %v, want %v", got, tt.taskCount)
			}
			if got := peak.Load(); got > int64(tt.n) {
				t.Errorf("RunAll() peak concurrency = %v, want <= %v", got, tt.n)
			}
		})
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
			semaphorex.RunAll(n, nil)
		})
	}
}

func BenchmarkRunAll(b *testing.B) {
	tasks := make([]func(), 1000)
	for i := range tasks {
		tasks[i] = func() {}
	}
	for b.Loop() {
		semaphorex.RunAll(10, tasks)
	}
}
