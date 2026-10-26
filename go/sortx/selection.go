package sortx

import "cmp"

func SelectionSort[S ~[]E, E cmp.Ordered](s S) {
	for i := range len(s) - 1 {
		minIdx := i
		for j := i + 1; j < len(s); j++ {
			if cmp.Less(s[j], s[minIdx]) {
				minIdx = j
			}
		}
		if minIdx != i {
			s[i], s[minIdx] = s[minIdx], s[i]
		}
	}
}
