package mutex_test

import (
	"sync"
	"testing"

	"github.com/palebluedot4/quark/go/concurrency/mutex"
)

func TestOpportunisticRunner(t *testing.T) {
	t.Parallel()
	t.Run("contended", func(t *testing.T) {
		t.Parallel()
		var r mutex.OpportunisticRunner
		started := make(chan struct{})
		release := make(chan struct{})
		var wg sync.WaitGroup
		wg.Go(func() {
			if got := r.TryRun(func() {
				close(started)
				<-release
			}); !got {
				t.Errorf("TryRun() = %v, want %v", got, true)
				close(started)
			}
		})
		<-started
		if got := r.TryRun(func() {}); got {
			t.Errorf("TryRun() = %v, want %v", got, false)
		}
		if got := r.TryRun(func() {}); got {
			t.Errorf("TryRun() = %v, want %v", got, false)
		}
		close(release)
		wg.Wait()
		if got := r.Runs(); got != 1 {
			t.Errorf("Runs() = %v, want %v", got, 1)
		}
		if got := r.Skips(); got != 2 {
			t.Errorf("Skips() = %v, want %v", got, 2)
		}
	})

	t.Run("uncontended", func(t *testing.T) {
		t.Parallel()
		var (
			r   mutex.OpportunisticRunner
			ran int
		)
		if got := r.TryRun(func() { ran++ }); !got {
			t.Errorf("TryRun() = %v, want %v", got, true)
		}
		if got := r.TryRun(func() { ran++ }); !got {
			t.Errorf("TryRun() = %v, want %v", got, true)
		}
		if ran != 2 {
			t.Errorf("work executions = %v, want %v", ran, 2)
		}
		if got := r.Runs(); got != 2 {
			t.Errorf("Runs() = %v, want %v", got, 2)
		}
		if got := r.Skips(); got != 0 {
			t.Errorf("Skips() = %v, want %v", got, 0)
		}
	})
}
