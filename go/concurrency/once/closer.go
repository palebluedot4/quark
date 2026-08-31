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
