package pool_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/palebluedot4/quark/go/concurrency/pool"
)

func TestGet(t *testing.T) {
	t.Parallel()
	buf := pool.Get()
	if buf == nil {
		t.Fatal("Get() = nil, want non-nil buffer")
	}
	if got := buf.Len(); got != 0 {
		t.Errorf("Get().Len() = %d, want 0", got)
	}
}

func TestPut(t *testing.T) {
	tests := []struct {
		name string
		grow int
		want bool
	}{
		{
			name: "normal buffer",
			grow: 8 << 10,
			want: true,
		},
		{
			name: "oversized buffer",
			grow: 512 << 10,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			buf.Grow(tt.grow)
			buf.WriteString("payload")
			pool.Put(buf)
			got := buf.Len() == 0
			if got != tt.want {
				t.Errorf("Put() reset buffer = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrent(t *testing.T) {
	t.Parallel()
	const (
		workers    = 100
		iterations = 1000
	)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for range iterations {
				buf := pool.Get()
				if got := buf.Len(); got != 0 {
					t.Errorf("Get().Len() = %d, want 0", got)
				}
				buf.WriteString("concurrent")
				pool.Put(buf)
			}
		})
	}
	wg.Wait()
}
