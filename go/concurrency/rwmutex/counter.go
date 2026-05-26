package rwmutex

import "sync"

type Counter struct {
	mu    sync.RWMutex
	value uint64
}

func (c *Counter) Value() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Add(delta uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

func (c *Counter) Set(val uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = val
}
