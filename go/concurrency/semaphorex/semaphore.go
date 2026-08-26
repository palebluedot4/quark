package semaphorex

import "context"

type Semaphore struct {
	tokens chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	if n < 1 {
		panic("semaphorex.NewSemaphore: limit must be positive")
	}
	return &Semaphore{tokens: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.tokens <- struct{}{}:
		return nil
	}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.tokens <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release() {
	select {
	case <-s.tokens:
	default:
		panic("semaphorex.Semaphore.Release: release of unacquired semaphore")
	}
}
