package waitgroup_test

import (
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/palebluedot4/quark/go/concurrency/waitgroup"
)

func TestRunAll(t *testing.T) {
	t.Parallel()
	runners := []struct {
		name string
		f    func([]func())
	}{
		{name: "RunAll", f: waitgroup.RunAll},
		{name: "RunAllManual", f: waitgroup.RunAllManual},
	}
	tests := []struct {
		name  string
		count int
	}{
		{name: "empty tasks", count: 0},
		{name: "single task", count: 1},
		{name: "concurrent 100", count: 100},
		{name: "concurrent 10000", count: 10000},
	}

	for _, runner := range runners {
		t.Run(runner.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					var ran atomic.Uint64
					tasks := make([]func(), tt.count)
					for i := range tasks {
						tasks[i] = func() { ran.Add(1) }
					}
					runner.f(tasks)
					got := ran.Load()
					if want := uint64(tt.count); got != want {
						t.Errorf("%s() executed tasks = %v, want %v", runner.name, got, want)
					}
				})
			}
		})
	}
}

func TestRunAllStartsTasksConcurrently(t *testing.T) {
	t.Parallel()
	runners := []struct {
		name string
		f    func([]func())
	}{
		{name: "RunAll", f: waitgroup.RunAll},
		{name: "RunAllManual", f: waitgroup.RunAllManual},
	}

	for _, runner := range runners {
		t.Run(runner.name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				started := make(chan struct{}, 2)
				release := make(chan struct{})
				tasks := []func(){
					func() {
						started <- struct{}{}
						<-release
					},
					func() {
						started <- struct{}{}
						<-release
					},
				}
				done := make(chan struct{})
				go func() {
					runner.f(tasks)
					close(done)
				}()
				synctest.Wait()
				if got := len(started); got != len(tasks) {
					t.Errorf("%s() concurrently started tasks = %v, want %v", runner.name, got, len(tasks))
				}
				close(release)
				<-done
			})
		})
	}
}

func BenchmarkRunAll(b *testing.B) {
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
