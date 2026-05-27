package mutex

import (
	"sync"
	"sync/atomic"
)

type OpportunisticRunner struct {
	mu    sync.Mutex
	runs  atomic.Uint64
	skips atomic.Uint64
}

func (r *OpportunisticRunner) TryRun(work func()) bool {
	if !r.mu.TryLock() {
		r.skips.Add(1)
		return false
	}
	defer r.mu.Unlock()
	r.runs.Add(1)
	work()
	return true
}

func (r *OpportunisticRunner) Runs() uint64 {
	return r.runs.Load()
}

func (r *OpportunisticRunner) Skips() uint64 {
	return r.skips.Load()
}
