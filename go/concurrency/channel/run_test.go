package channel_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/channel"
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
			synctest.Test(t, func(t *testing.T) {
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
						time.Sleep(time.Second)
						current.Add(-1)
					}
				}
				channel.RunAll(tt.n, tasks)
				if got := ran.Load(); got != int64(tt.taskCount) {
					t.Errorf("RunAll() executed tasks = %v, want %v", got, tt.taskCount)
				}
				got := peak.Load()
				want := int64(min(tt.n, tt.taskCount))
				if got != want {
					t.Errorf("RunAll() peak concurrency = %v, want %v", got, want)
				}
			})
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
			channel.RunAll(n, nil)
		})
	}
}

func BenchmarkRunAll(b *testing.B) {
	tasks := make([]func(), 1000)
	for i := range tasks {
		tasks[i] = func() {}
	}
	for b.Loop() {
		channel.RunAll(10, tasks)
	}
}
