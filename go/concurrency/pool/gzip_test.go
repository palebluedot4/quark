package pool_test

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/pool"
)

func TestGetGzipWriter(t *testing.T) {
	t.Parallel()
	zw := pool.GetGzipWriter(io.Discard)
	zw.Name = "stale.txt"
	if err := pool.PutGzipWriter(zw); err != nil {
		t.Fatalf("PutGzipWriter() = %v, want nil", err)
	}
	zw = pool.GetGzipWriter(io.Discard)
	if got := zw.Name; got != "" {
		t.Errorf("GetGzipWriter().Name = %q, want %q", got, "")
	}
	if err := pool.PutGzipWriter(zw); err != nil {
		t.Fatalf("PutGzipWriter() = %v, want nil", err)
	}
}

func TestPutGzipWriter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
	}{
		{name: "empty", in: ""},
		{name: "payload", in: "payload"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			zw := pool.GetGzipWriter(&buf)
			if _, err := io.WriteString(zw, tt.in); err != nil {
				t.Fatalf("Write(%q) = _, %v, want nil", tt.in, err)
			}
			if err := pool.PutGzipWriter(zw); err != nil {
				t.Fatalf("PutGzipWriter() = %v, want nil", err)
			}
			got, err := gunzip(&buf)
			if err != nil || got != tt.in {
				t.Errorf("gunzip() = %q, %v, want %q, nil", got, err, tt.in)
			}
		})
	}
}

func gunzip(r io.Reader) (string, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(zr)
	return string(b), err
}

func TestPutGzipWriterReturnsCloseError(t *testing.T) {
	t.Parallel()
	want := errors.New("write failed")
	zw := pool.GetGzipWriter(failingWriter{err: want})
	if got := pool.PutGzipWriter(zw); !errors.Is(got, want) {
		t.Errorf("PutGzipWriter() = %v, want %v", got, want)
	}
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

//nolint:paralleltest
func TestPutGzipWriterReleasesDestination(t *testing.T) {
	// This test must not run in parallel with other tests as it asserts on when
	// the destination becomes unreachable, and their allocations delay the
	// cleanup queue it waits on.
	var released atomic.Bool
	func() {
		buf := new(bytes.Buffer)
		runtime.AddCleanup(buf, func(b *atomic.Bool) { b.Store(true) }, &released)
		if err := pool.PutGzipWriter(pool.GetGzipWriter(buf)); err != nil {
			t.Fatalf("PutGzipWriter() = %v, want nil", err)
		}
	}()
	runtime.GC()
	if !cleanedUp(&released) {
		t.Error("PutGzipWriter() retained the destination, want released")
	}
}

func cleanedUp(b *atomic.Bool) bool {
	for range 100 {
		if b.Load() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}

func TestGzipWriterConcurrent(t *testing.T) {
	t.Parallel()
	const (
		workers    = 10
		iterations = 100
		payload    = "concurrent"
	)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for range iterations {
				var buf bytes.Buffer
				zw := pool.GetGzipWriter(&buf)
				if _, err := io.WriteString(zw, payload); err != nil {
					t.Errorf("Write(%q) = _, %v, want nil", payload, err)
				}
				if err := pool.PutGzipWriter(zw); err != nil {
					t.Errorf("PutGzipWriter() = %v, want nil", err)
				}
				if got, err := gunzip(&buf); err != nil || got != payload {
					t.Errorf("gunzip() = %q, %v, want %q, nil", got, err, payload)
				}
			}
		})
	}
	wg.Wait()
}
