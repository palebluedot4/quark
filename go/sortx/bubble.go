package sortx

import "cmp"

func BubbleSort[S ~[]E, E cmp.Ordered](s S) {
	n := len(s)
	for n > 1 {
		lastSwap := 0
		for j := range n - 1 {
			if cmp.Less(s[j+1], s[j]) {
				s[j], s[j+1] = s[j+1], s[j]
				lastSwap = j + 1
			}
		}
		n = lastSwap
	}
}
