package pool

import (
	"bytes"
	"sync"
)

const maxBufferSize = 64 << 10

var buffers = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	if buf.Cap() > maxBufferSize {
		return
	}
	buf.Reset()
	buffers.Put(buf)
}
