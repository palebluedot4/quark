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

type contender[T any] struct {
	value    T
	defeated []T
}

func TopTwoTournament[S ~[]E, E cmp.Ordered](s S) (first, second E) {
	if len(s) < 2 {
		panic("selection.TopTwoTournament: slice too short")
	}
	contenders := make([]contender[E], len(s))
	for i, v := range s {
		contenders[i] = contender[E]{value: v}
	}
	for len(contenders) > 1 {
		next := make([]contender[E], 0, (len(contenders)+1)/2)
		for i := 0; i+1 < len(contenders); i += 2 {
			winner, loser := i, i+1
			if cmp.Less(contenders[winner].value, contenders[loser].value) {
				winner, loser = loser, winner
			}
			advancing := contenders[winner]
			advancing.defeated = append(advancing.defeated, contenders[loser].value)
			next = append(next, advancing)
		}
		if len(contenders)%2 == 1 {
			next = append(next, contenders[len(contenders)-1])
		}
		contenders = next
	}
	champion := contenders[0]
	first, second = champion.value, champion.defeated[0]
	for _, candidate := range champion.defeated[1:] {
		if cmp.Less(second, candidate) {
			second = candidate
		}
	}
	return first, second
}
