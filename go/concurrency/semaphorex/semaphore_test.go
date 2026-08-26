package semaphorex_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/palebluedot4/quark/go/concurrency/semaphorex"
)

func TestSemaphoreAcquire(t *testing.T) {
	t.Parallel()
	t.Run("blocks at limit", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			sem := semaphorex.NewSemaphore(1)
			if err := sem.Acquire(context.Background()); err != nil {
				t.Fatalf("Acquire() = %v, want nil", err)
			}
			done := make(chan error, 1)
			go func() {
				done <- sem.Acquire(context.Background())
			}()
			synctest.Wait()
			select {
			case err := <-done:
				t.Fatalf("Acquire() = %v before Release, want it to block", err)
			default:
			}
			sem.Release()
			if err := <-done; err != nil {
				t.Errorf("Acquire() = %v, want nil", err)
			}
			sem.Release()
		})
	})

	t.Run("canceled while blocked", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			sem := semaphorex.NewSemaphore(1)
			if err := sem.Acquire(context.Background()); err != nil {
				t.Fatalf("Acquire() = %v, want nil", err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() {
				done <- sem.Acquire(ctx)
			}()
			synctest.Wait()
			cancel()
			if err := <-done; !errors.Is(err, context.Canceled) {
				t.Errorf("Acquire() = %v, want %v", err, context.Canceled)
			}
			sem.Release()
			if got := sem.TryAcquire(); !got {
				t.Errorf("TryAcquire() = %v, want %v", got, true)
			}
			sem.Release()
		})
	})

	t.Run("already canceled", func(t *testing.T) {
		t.Parallel()
		sem := semaphorex.NewSemaphore(1)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := sem.Acquire(ctx); !errors.Is(err, context.Canceled) {
			t.Errorf("Acquire() = %v, want %v", err, context.Canceled)
		}
		if got := sem.TryAcquire(); !got {
			t.Errorf("TryAcquire() = %v, want %v", got, true)
		}
		sem.Release()
	})
}

func TestSemaphoreTryAcquire(t *testing.T) {
	t.Parallel()
	sem := semaphorex.NewSemaphore(1)
	if got := sem.TryAcquire(); !got {
		t.Errorf("TryAcquire() = %v, want %v", got, true)
	}
	if got := sem.TryAcquire(); got {
		t.Errorf("TryAcquire() = %v, want %v", got, false)
	}
	sem.Release()
	if got := sem.TryAcquire(); !got {
		t.Errorf("TryAcquire() = %v, want %v", got, true)
	}
	sem.Release()
}

func TestSemaphoreConcurrent(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		const (
			limit = 5
			tasks = 100
		)
		sem := semaphorex.NewSemaphore(limit)
		var (
			active atomic.Int64
			peak   atomic.Int64
			wg     sync.WaitGroup
		)
		for range tasks {
			wg.Go(func() {
				if err := sem.Acquire(context.Background()); err != nil {
					t.Errorf("Acquire() = %v, want nil", err)
					return
				}
				current := active.Add(1)
				for {
					previous := peak.Load()
					if current <= previous || peak.CompareAndSwap(previous, current) {
						break
					}
				}
				time.Sleep(time.Second)
				active.Add(-1)
				sem.Release()
			})
		}
		wg.Wait()
		if got := peak.Load(); got != limit {
			t.Errorf("peak concurrency = %v, want %v", got, limit)
		}
	})
}

func TestNewSemaphorePanicsOnNonPositiveLimit(t *testing.T) {
	t.Parallel()
	for _, limit := range []int{0, -1} {
		t.Run(fmt.Sprintf("limit=%d", limit), func(t *testing.T) {
			t.Parallel()
			defer func() {
				if recover() == nil {
					t.Errorf("NewSemaphore(%d) did not panic, want panic", limit)
				}
			}()
			semaphorex.NewSemaphore(limit)
		})
	}
}

func TestSemaphoreReleasePanicsOnUnacquired(t *testing.T) {
	t.Parallel()
	sem := semaphorex.NewSemaphore(1)
	defer func() {
		if recover() == nil {
			t.Error("Release() did not panic, want panic")
		}
	}()
	sem.Release()
}
