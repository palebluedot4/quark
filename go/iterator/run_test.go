package iterator_test

import (
	"errors"
	"iter"
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/palebluedot4/quark/go/iterator"
)

func TestRunAll(t *testing.T) {
	t.Parallel()
	runners := []struct {
		name string
		f    func([]func())
	}{
		{name: "RunAll", f: func(tasks []func()) { iterator.RunAll(slices.Values(tasks)) }},
		{name: "RunAllManual", f: func(tasks []func()) { iterator.RunAllManual(slices.Values(tasks)) }},
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

func TestRunAllStartsTasksConcurrently(t *testing.T) {
	t.Parallel()
	runners := []struct {
		name string
		f    func([]func())
	}{
		{name: "RunAll", f: func(tasks []func()) { iterator.RunAll(slices.Values(tasks)) }},
		{name: "RunAllManual", f: func(tasks []func()) { iterator.RunAllManual(slices.Values(tasks)) }},
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

func TestRunUntilError(t *testing.T) {
	t.Parallel()
	runners := []struct {
		name string
		f    func(iter.Seq[func() error]) error
	}{
		{name: "RunUntilError", f: iterator.RunUntilError},
		{name: "RunUntilErrorManual", f: iterator.RunUntilErrorManual},
	}
	want := errors.New("boom")

	for _, runner := range runners {
		t.Run(runner.name, func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				name      string
				taskCount int
				failAt    int
				want      error
				wantRan   int
			}{
				{
					name:      "empty tasks",
					taskCount: 0,
					failAt:    -1,
					want:      nil,
					wantRan:   0,
				},
				{
					name:      "no errors",
					taskCount: 3,
					failAt:    -1,
					want:      nil,
					wantRan:   3,
				},
				{
					name:      "fails on first task",
					taskCount: 3,
					failAt:    0,
					want:      want,
					wantRan:   1,
				},
				{
					name:      "fails on middle task",
					taskCount: 3,
					failAt:    1,
					want:      want,
					wantRan:   2,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					var ran int
					tasks := make([]func() error, tt.taskCount)
					for i := range tasks {
						tasks[i] = func() error {
							ran++
							if i == tt.failAt {
								return want
							}
							return nil
						}
					}
					got := runner.f(slices.Values(tasks))
					if !errors.Is(got, tt.want) {
						t.Errorf("%s() = %v, want %v", runner.name, got, tt.want)
					}
					if ran != tt.wantRan {
						t.Errorf("%s() ran %v tasks, want %v", runner.name, ran, tt.wantRan)
					}
				})
			}

			t.Run("stops source iteration after error", func(t *testing.T) {
				t.Parallel()
				yielded := 0
				tasks := func(yield func(func() error) bool) {
					for i := range 3 {
						yielded++
						task := func() error {
							if i == 1 {
								return want
							}
							return nil
						}
						if !yield(task) {
							return
						}
					}
				}
				if err := runner.f(tasks); !errors.Is(err, want) {
					t.Errorf("%s() = %v, want %v", runner.name, err, want)
				}
				if yielded != 2 {
					t.Errorf("%s() consumed tasks = %v, want %v", runner.name, yielded, 2)
				}
			})
		})
	}
}
