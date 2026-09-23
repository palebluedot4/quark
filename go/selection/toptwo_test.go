package selection_test

import (
	"cmp"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/palebluedot4/quark/go/selection"
)

type topTwoVariant[S ~[]E, E cmp.Ordered] struct {
	name string
	f    func(S) (E, E)
}

func topTwoVariants[S ~[]E, E cmp.Ordered]() []topTwoVariant[S, E] {
	return []topTwoVariant[S, E]{
		{name: "TopTwoLinear", f: selection.TopTwoLinear[S]},
		{name: "TopTwoTournament", f: selection.TopTwoTournament[S]},
		{name: "TopTwoSorted", f: selection.TopTwoSorted[S]},
	}
}

func TestTopTwo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		in         []int
		wantFirst  int
		wantSecond int
	}{
		{
			name:       "ascending pair",
			in:         []int{1, 2},
			wantFirst:  2,
			wantSecond: 1,
		},
		{
			name:       "descending pair",
			in:         []int{2, 1},
			wantFirst:  2,
			wantSecond: 1,
		},
		{
			name:       "maximum last",
			in:         []int{1, 2, 3, 4, 5},
			wantFirst:  5,
			wantSecond: 4,
		},
		{
			name:       "maximum first",
			in:         []int{5, 4, 3, 2, 1},
			wantFirst:  5,
			wantSecond: 4,
		},
		{
			name:       "runner-up before maximum",
			in:         []int{2, 3, 1},
			wantFirst:  3,
			wantSecond: 2,
		},
		{
			name:       "duplicate maximum",
			in:         []int{5, 3, 5},
			wantFirst:  5,
			wantSecond: 5,
		},
		{
			name:       "negative",
			in:         []int{-3, -1, -2},
			wantFirst:  -1,
			wantSecond: -2,
		},
	}

	for _, v := range topTwoVariants[[]int]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					first, second := v.f(tt.in)
					if first != tt.wantFirst || second != tt.wantSecond {
						t.Errorf("%s(%v) = %v, %v, want %v, %v", v.name, tt.in, first, second, tt.wantFirst, tt.wantSecond)
					}
				})
			}
		})
	}
}

func TestTopTwoSupportsOrderedTypes(t *testing.T) {
	t.Parallel()
	type scores []int
	testTopTwo(t, scores{1, 3, 2}, 3, 2)
	testTopTwo(t, []string{"a", "c", "b"}, "c", "b")
}

func testTopTwo[S ~[]E, E cmp.Ordered](t *testing.T, in S, wantFirst, wantSecond E) {
	t.Helper()
	t.Run(reflect.TypeFor[S]().String(), func(t *testing.T) {
		t.Parallel()
		for _, v := range topTwoVariants[S]() {
			t.Run(v.name, func(t *testing.T) {
				t.Parallel()
				first, second := v.f(in)
				if first != wantFirst || second != wantSecond {
					t.Errorf("%s(%v) = %v, %v, want %v, %v", v.name, in, first, second, wantFirst, wantSecond)
				}
			})
		}
	})
}

func TestTopTwoOrdersNaNLast(t *testing.T) {
	t.Parallel()
	nan := math.NaN()
	tests := []struct {
		name       string
		in         []float64
		wantFirst  float64
		wantSecond float64
	}{
		{
			name:       "leading NaN",
			in:         []float64{nan, 1},
			wantFirst:  1,
			wantSecond: nan,
		},
		{
			name:       "trailing NaN",
			in:         []float64{1, nan},
			wantFirst:  1,
			wantSecond: nan,
		},
		{
			name:       "NaN between numbers",
			in:         []float64{1, nan, 2},
			wantFirst:  2,
			wantSecond: 1,
		},
		{
			name:       "NaN after numbers",
			in:         []float64{1, 2, nan},
			wantFirst:  2,
			wantSecond: 1,
		},
		{
			name:       "single non-NaN",
			in:         []float64{nan, nan, 1},
			wantFirst:  1,
			wantSecond: nan,
		},
	}

	for _, v := range topTwoVariants[[]float64]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					first, second := v.f(tt.in)
					if !equalFloat(first, tt.wantFirst) || !equalFloat(second, tt.wantSecond) {
						t.Errorf("%s(%v) = %v, %v, want %v, %v", v.name, tt.in, first, second, tt.wantFirst, tt.wantSecond)
					}
				})
			}
		})
	}
}

func equalFloat(got, want float64) bool {
	if math.IsNaN(want) {
		return math.IsNaN(got)
	}
	return got == want
}

func TestTopTwoDoesNotModifyInput(t *testing.T) {
	t.Parallel()
	for _, v := range topTwoVariants[[]int]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			in := []int{3, 1, 4, 1, 5}
			want := slices.Clone(in)
			v.f(in)
			if !slices.Equal(in, want) {
				t.Errorf("after %s(%v), input = %v, want %v", v.name, want, in, want)
			}
		})
	}
}

func TestTopTwoPanicsOnShortSlice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []int
	}{
		{name: "nil", in: nil},
		{name: "empty", in: []int{}},
		{name: "single", in: []int{1}},
	}

	for _, v := range topTwoVariants[[]int]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					defer func() {
						if recover() == nil {
							t.Errorf("%s(%#v) did not panic, want panic", v.name, tt.in)
						}
					}()
					v.f(tt.in)
				})
			}
		})
	}
}
