package channel_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/palebluedot4/quark/go/concurrency/channel"
)

func TestQueue(t *testing.T) {
	t.Parallel()
	t.Run("fifo order", func(t *testing.T) {
		t.Parallel()
		q := channel.NewQueue[int](3)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		for _, want := range []int{1, 2, 3} {
			if got := q.Dequeue(); got != want {
				t.Errorf("Dequeue() = %v, want %v", got, want)
			}
		}
	})

	t.Run("nonblocking enqueue", func(t *testing.T) {
		t.Parallel()
		q := channel.NewQueue[int](1)
		if got := q.TryEnqueue(1); !got {
			t.Errorf("TryEnqueue(1) = %v, want %v", got, true)
		}
		if got := q.TryEnqueue(2); got {
			t.Errorf("TryEnqueue(2) = %v, want %v", got, false)
		}
		if got := q.Dequeue(); got != 1 {
			t.Errorf("Dequeue() = %v, want %v", got, 1)
		}
		if got := q.TryEnqueue(2); !got {
			t.Errorf("TryEnqueue(2) = %v, want %v", got, true)
		}
		if got := q.Dequeue(); got != 2 {
			t.Errorf("Dequeue() = %v, want %v", got, 2)
		}
	})

	t.Run("nonblocking dequeue", func(t *testing.T) {
		t.Parallel()
		q := channel.NewQueue[int](1)
		if got, ok := q.TryDequeue(); ok || got != 0 {
			t.Errorf("TryDequeue() = %v, %v, want %v, %v", got, ok, 0, false)
		}
		q.Enqueue(1)
		if got, ok := q.TryDequeue(); !ok || got != 1 {
			t.Errorf("TryDequeue() = %v, %v, want %v, %v", got, ok, 1, true)
		}
		if got, ok := q.TryDequeue(); ok || got != 0 {
			t.Errorf("TryDequeue() = %v, %v, want %v, %v", got, ok, 0, false)
		}
	})

	t.Run("blocks on empty", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			q := channel.NewQueue[int](1)
			done := make(chan int, 1)
			go func() {
				done <- q.Dequeue()
			}()
			synctest.Wait()
			select {
			case got := <-done:
				t.Fatalf("Dequeue() = %v before Enqueue, want it to block", got)
			default:
			}
			q.Enqueue(42)
			if got := <-done; got != 42 {
				t.Errorf("Dequeue() = %v, want %v", got, 42)
			}
		})
	})

	t.Run("blocks on full", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			q := channel.NewQueue[int](1)
			q.Enqueue(1)
			done := make(chan struct{}, 1)
			go func() {
				q.Enqueue(2)
				done <- struct{}{}
			}()
			synctest.Wait()
			select {
			case <-done:
				t.Fatal("Enqueue() returned before space freed, want it to block")
			default:
			}
			if got := q.Dequeue(); got != 1 {
				t.Errorf("Dequeue() = %v, want %v", got, 1)
			}
			<-done
			if got := q.Dequeue(); got != 2 {
				t.Errorf("Dequeue() = %v, want %v", got, 2)
			}
		})
	})

	t.Run("concurrent producers and consumers", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			const (
				producers        = 4
				consumers        = 3
				itemsPerProducer = 25
				total            = producers * itemsPerProducer
			)
			q := channel.NewQueue[int](7)
			items := make(chan int, total)
			var (
				next atomic.Int64
				wg   sync.WaitGroup
			)
			for range consumers {
				wg.Go(func() {
					for next.Add(1) <= total {
						items <- q.Dequeue()
					}
				})
			}
			for producer := range producers {
				wg.Go(func() {
					for item := range itemsPerProducer {
						q.Enqueue(producer*itemsPerProducer + item)
					}
				})
			}
			wg.Wait()
			close(items)
			counts := make([]int, total)
			for item := range items {
				if item < 0 || item >= total {
					t.Errorf("Dequeue() = %v, want value in [0, %d)", item, total)
					continue
				}
				counts[item]++
			}
			for item, count := range counts {
				if count != 1 {
					t.Errorf("Dequeue() count for value %d = %v, want %v", item, count, 1)
				}
			}
		})
	})

	t.Run("invalid capacity", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			capacity int
		}{
			{name: "zero", capacity: 0},
			{name: "negative", capacity: -1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewQueue(%d) did not panic, want panic", tt.capacity)
					}
				}()
				channel.NewQueue[int](tt.capacity)
			})
		}
	})
}
