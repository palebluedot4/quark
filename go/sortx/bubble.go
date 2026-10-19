package sortx

import "cmp"

func BubbleSort[S ~[]E, E cmp.Ordered](s S) {
	n := len(s)
	for n > 1 {
		lastSwap := 0
		for i := range n - 1 {
			if cmp.Less(s[i+1], s[i]) {
				s[i], s[i+1] = s[i+1], s[i]
				lastSwap = i + 1
			}
		}
		n = lastSwap
	}
}

func BubbleSortFlag[S ~[]E, E cmp.Ordered](s S) {
	for i := range len(s) - 1 {
		swapped := false
		for j := range len(s) - i - 1 {
			if cmp.Less(s[j+1], s[j]) {
				s[j], s[j+1] = s[j+1], s[j]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}
