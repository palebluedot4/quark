package unsafes

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Once struct {
	sync.Once
}

func (o *Once) Done() bool {
	// This relies on the internal memory layout of sync.Once.
	// This is unsafe and may break if the standard library implementation changes.
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&o.Once))) == 1
}
