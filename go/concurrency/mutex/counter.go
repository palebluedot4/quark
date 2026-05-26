package mutex

import "sync"

type Counter struct {
	mu    sync.Mutex
	value uint64
}

func (c *Counter) Value() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
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
