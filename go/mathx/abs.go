package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

func AbsSigned[T constraints.Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func AbsFloat[T constraints.Float](x T) T {
	return T(math.Abs(float64(x)))
}
