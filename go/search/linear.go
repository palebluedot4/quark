package search

import "slices"

func LinearSearch[S ~[]E, E comparable](s S, target E) int {
	return slices.Index(s, target)
}

func LinearSearchManual[S ~[]E, E comparable](s S, target E) int {
	for i := range s {
		if s[i] == target {
			return i
		}
	}
	return -1
}
