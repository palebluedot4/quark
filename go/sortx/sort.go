package sortx

import (
	"cmp"
	"slices"
)

func Sort[S ~[]E, E cmp.Ordered](s S) {
	slices.Sort(s)
}
