package rwmutex_test

import (
	"sync"
	"testing"

	"github.com/palebluedot4/quark/go/concurrency/rwmutex"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		fn   func(*rwmutex.Counter)
		want uint64
	}{
		{
			name: "increment",
			fn:   func(c *rwmutex.Counter) { c.Increment() },
			want: 1,
		},
		{
			name: "add",
			fn:   func(c *rwmutex.Counter) { c.Add(5) },
			want: 5,
		},
		{
			name: "set",
			fn:   func(c *rwmutex.Counter) { c.Set(42) },
			want: 42,
		},
		{
			name: "combined",
			fn: func(c *rwmutex.Counter) {
				c.Increment()
				c.Add(10)
				c.Set(5)
				c.Add(2)
			},
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var c rwmutex.Counter
			tt.fn(&c)
			got := c.Value()
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounter_Concurrent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		fn         func(*rwmutex.Counter)
		workers    int
		iterations int
		want       uint64
	}{
		{
			name:       "increment",
			fn:         func(c *rwmutex.Counter) { c.Increment() },
			workers:    100,
			iterations: 1000,
			want:       100 * 1000,
		},
		{
			name:       "add",
			fn:         func(c *rwmutex.Counter) { c.Add(5) },
			workers:    50,
			iterations: 1000,
			want:       50 * 1000 * 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				c  rwmutex.Counter
				wg sync.WaitGroup
			)
			for range tt.workers {
				wg.Go(func() {
					for range tt.iterations {
						tt.fn(&c)
					}
				})
			}
			wg.Wait()
			got := c.Value()
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCounter_Value(b *testing.B) {
	var c rwmutex.Counter
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Value()
		}
	})
}

func BenchmarkCounter_Increment(b *testing.B) {
	var c rwmutex.Counter
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Increment()
		}
	})
}
