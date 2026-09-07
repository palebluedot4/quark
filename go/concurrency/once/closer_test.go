package once_test

import (
	"errors"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/once"
)

var errCleanup = errors.New("cleanup failed")

func TestCloser(t *testing.T) {
	t.Parallel()
	variants := []struct {
		name string
		f    func(func() error) io.Closer
	}{
		{name: "Closer", f: func(cleanup func() error) io.Closer { return once.NewCloser(cleanup) }},
		{name: "CloserManual", f: func(cleanup func() error) io.Closer { return once.NewCloserManual(cleanup) }},
	}
	tests := []struct {
		name string
		want error
	}{
		{name: "success", want: nil},
		{name: "failure", want: errCleanup},
	}

	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					calls := 0
					c := v.f(func() error {
						calls++
						return tt.want
					})
					for range 3 {
						if got := c.Close(); !errors.Is(got, tt.want) {
							t.Errorf("%s.Close() = %v, want %v", v.name, got, tt.want)
						}
					}
					if calls != 1 {
						t.Errorf("%s cleanup calls = %v, want %v", v.name, calls, 1)
					}
				})
			}
		})
	}
}

func TestCloserConcurrent(t *testing.T) {
	t.Parallel()
	variants := []struct {
		name string
		f    func(func() error) io.Closer
	}{
		{name: "Closer", f: func(cleanup func() error) io.Closer { return once.NewCloser(cleanup) }},
		{name: "CloserManual", f: func(cleanup func() error) io.Closer { return once.NewCloserManual(cleanup) }},
	}

	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			const goroutines = 100
			var calls atomic.Int64
			c := v.f(func() error {
				calls.Add(1)
				return errCleanup
			})
			errs := make([]error, goroutines)
			var wg sync.WaitGroup
			for i := range errs {
				wg.Go(func() {
					errs[i] = c.Close()
				})
			}
			wg.Wait()
			if got := calls.Load(); got != 1 {
				t.Errorf("%s cleanup calls = %v, want %v", v.name, got, 1)
			}
			for _, err := range errs {
				if !errors.Is(err, errCleanup) {
					t.Errorf("%s.Close() = %v, want %v", v.name, err, errCleanup)
				}
			}
		})
	}
}

func TestCloserPanic(t *testing.T) {
	t.Parallel()
	t.Run("Closer", func(t *testing.T) {
		t.Parallel()
		calls := 0
		c := once.NewCloser(func() error {
			calls++
			panic("boom")
		})
		for range 2 {
			func() {
				defer func() {
					if p := recover(); p != "boom" {
						t.Errorf("Closer.Close() panicked with %v, want %v", p, "boom")
					}
				}()
				c.Close()
				t.Error("Closer.Close() returned, want panic")
			}()
		}
		if calls != 1 {
			t.Errorf("Closer cleanup calls = %v, want %v", calls, 1)
		}
	})

	t.Run("CloserManual", func(t *testing.T) {
		t.Parallel()
		calls := 0
		c := once.NewCloserManual(func() error {
			calls++
			panic("boom")
		})
		func() {
			defer func() {
				if p := recover(); p != "boom" {
					t.Errorf("CloserManual.Close() panicked with %v, want %v", p, "boom")
				}
			}()
			c.Close()
			t.Error("CloserManual.Close() returned, want panic")
		}()
		if got := c.Close(); got != nil {
			t.Errorf("CloserManual.Close() = %v, want %v", got, error(nil))
		}
		if calls != 1 {
			t.Errorf("CloserManual cleanup calls = %v, want %v", calls, 1)
		}
	})
}

//nolint:paralleltest
func TestCloserReleasesCleanup(t *testing.T) {
	// This test must not run in parallel with other tests as it asserts on when
	// the captured buffer becomes unreachable, and their allocations delay the
	// cleanup queue it drains with runtime.GC.
	variants := []struct {
		name string
		f    func(func() error) io.Closer
	}{
		{name: "Closer", f: func(cleanup func() error) io.Closer { return once.NewCloser(cleanup) }},
		{name: "CloserManual", f: func(cleanup func() error) io.Closer { return once.NewCloserManual(cleanup) }},
	}

	for _, v := range variants {
		t.Run(v.name, func(t *testing.T) {
			buf := make([]byte, 1024)
			var collected atomic.Bool
			runtime.AddCleanup(&buf[0], func(b *atomic.Bool) { b.Store(true) }, &collected)
			c := v.f(func() error {
				buf[0] = 1
				return nil
			})
			runtime.GC()
			if collected.Load() {
				t.Fatalf("%s collected the captured buffer before Close(), want retained", v.name)
			}
			if err := c.Close(); err != nil {
				t.Fatalf("%s.Close() = %v, want %v", v.name, err, error(nil))
			}
			if !collectedAfterGC(&collected) {
				t.Errorf("%s retained the captured buffer after Close(), want collected", v.name)
			}
			runtime.KeepAlive(c)
		})
	}
}

func collectedAfterGC(b *atomic.Bool) bool {
	for range 100 {
		runtime.GC()
		if b.Load() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}
