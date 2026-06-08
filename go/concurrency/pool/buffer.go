package pool

import (
	"bytes"
	"sync"
)

var buffers = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func Get() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}
