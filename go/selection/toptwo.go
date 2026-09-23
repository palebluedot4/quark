package selection

import "cmp"

func TopTwoLinear[S ~[]E, E cmp.Ordered](s S) (first, second E) {
	if len(s) < 2 {
		panic("selection.TopTwoLinear: slice too short")
	}
	first, second = s[0], s[1]
	if cmp.Less(first, second) {
		first, second = second, first
	}
	for i := 2; i < len(s); i++ {
		switch {
		case cmp.Less(first, s[i]):
			first, second = s[i], first
		case cmp.Less(second, s[i]):
			second = s[i]
		}
	}
	return first, second
}
