package mathx

import (
	"math"
	"unsafe"

	"golang.org/x/exp/constraints"
)

const bitsPerByte = 8

func AbsSigned[T constraints.Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func AbsSignedBitwise[T constraints.Signed](x T) T {
	mask := x >> (unsafe.Sizeof(x)*bitsPerByte - 1)
	return (x + mask) ^ mask
}

func AbsFloat[T constraints.Float](x T) T {
	return T(math.Abs(float64(x)))
}

func AbsFloatSignbit[T constraints.Float](x T) T {
	if math.Signbit(float64(x)) {
		return -x
	}
	return x
}

const (
	float32Size = 4
	float64Size = 8
)

func AbsFloatBitwise[T constraints.Float](x T) T {
	switch unsafe.Sizeof(x) {
	case float32Size:
		const signMask = 1 << (float32Size*bitsPerByte - 1)
		bits := *(*uint32)(unsafe.Pointer(&x)) &^ signMask
		return *(*T)(unsafe.Pointer(&bits))
	case float64Size:
		const signMask = 1 << (float64Size*bitsPerByte - 1)
		bits := *(*uint64)(unsafe.Pointer(&x)) &^ signMask
		return *(*T)(unsafe.Pointer(&bits))
	default:
		panic("mathx.AbsFloatBitwise: unreachable")
	}
}
