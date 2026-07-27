package pool_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/palebluedot4/quark/go/concurrency/pool"
)

func TestGetBuffer(t *testing.T) {
	t.Parallel()
	buf := pool.GetBuffer()
	if got := buf.Len(); got != 0 {
		t.Errorf("GetBuffer().Len() = %d, want 0", got)
	}
}

func TestPutBuffer(t *testing.T) {
	// This test must not run in parallel with other tests as PutBuffer returns
	// buf to the shared pool, where another test can take it and write to it.
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
			pool.PutBuffer(buf)
			got := buf.Len() == 0
			if got != tt.want {
				t.Errorf("PutBuffer() reset buffer = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBufferConcurrent(t *testing.T) {
	t.Parallel()
	const (
		workers    = 100
		iterations = 1000
	)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for range iterations {
				buf := pool.GetBuffer()
				if got := buf.Len(); got != 0 {
					t.Errorf("GetBuffer().Len() = %d, want 0", got)
				}
				buf.WriteString("concurrent")
				pool.PutBuffer(buf)
			}
		})
	}
	wg.Wait()
}

func BenchmarkBufferPooled(b *testing.B) {
	for b.Loop() {
		buf := pool.GetBuffer()
		buf.WriteString("benchmark payload")
		pool.PutBuffer(buf)
	}
}

func BenchmarkBufferUnpooled(b *testing.B) {
	for b.Loop() {
		buf := new(bytes.Buffer)
		buf.WriteString("benchmark payload")
	}
}
