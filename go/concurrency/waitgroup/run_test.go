package waitgroup_test

import (
	"sync/atomic"
	"testing"

	"github.com/palebluedot4/quark/go/concurrency/waitgroup"
)

func TestRun(t *testing.T) {
	t.Parallel()
	runners := []struct {
		name string
		f    func([]func())
	}{
		{name: "RunAll", f: waitgroup.RunAll},
		{name: "RunAllManual", f: waitgroup.RunAllManual},
	}
	tests := []struct {
		name string
		want uint64
	}{
		{name: "empty tasks", want: 0},
		{name: "single task", want: 1},
		{name: "concurrent 100", want: 100},
		{name: "concurrent 10000", want: 10000},
	}

	for _, runner := range runners {
		t.Run(runner.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					var counter atomic.Uint64
					tasks := make([]func(), tt.want)
					for i := range tasks {
						tasks[i] = func() { counter.Add(1) }
					}
					runner.f(tasks)
					got := counter.Load()
					if got != tt.want {
						t.Errorf("counter.Load() = %v, want %v", got, tt.want)
					}
				})
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	runners := []struct {
		name string
		f    func([]func())
	}{
		{name: "RunAll", f: waitgroup.RunAll},
		{name: "RunAllManual", f: waitgroup.RunAllManual},
	}
	tasks := make([]func(), 1000)
	for i := range tasks {
		tasks[i] = func() {}
	}

	for _, runner := range runners {
		b.Run(runner.name, func(b *testing.B) {
			for b.Loop() {
				runner.f(tasks)
			}
		})
	}
}
