package mathx

import "golang.org/x/exp/constraints"

func AbsSigned[T constraints.Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
