package once

import (
	"io"
	"sync"
)

type Closer struct {
	f func() error
}

var _ io.Closer = (*Closer)(nil)

func NewCloser(cleanup func() error) *Closer {
	return &Closer{f: sync.OnceValue(cleanup)}
}

func (c *Closer) Close() error {
	return c.f()
}

type CloserManual struct {
	once sync.Once
	f    func() error
	err  error
}

var _ io.Closer = (*CloserManual)(nil)

func NewCloserManual(cleanup func() error) *CloserManual {
	return &CloserManual{f: cleanup}
}

func (c *CloserManual) Close() error {
	c.once.Do(func() {
		defer func() { c.f = nil }()
		c.err = c.f()
	})
	return c.err
}
