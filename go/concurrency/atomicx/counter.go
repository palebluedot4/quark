package atomicx

import "sync/atomic"

type Counter struct {
	value atomic.Uint64
}

func (c *Counter) Value() uint64 {
	return c.value.Load()
}

func (c *Counter) Increment() {
	c.value.Add(1)
}

func (c *Counter) Add(delta uint64) {
	c.value.Add(delta)
}

func (c *Counter) Set(val uint64) {
	c.value.Store(val)
}
