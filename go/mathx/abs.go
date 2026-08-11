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
