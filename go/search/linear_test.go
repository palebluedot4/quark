package search_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/palebluedot4/quark/go/search"
)

type linearVariant[S ~[]E, E comparable] struct {
	name string
	f    func(S, E) int
}

func linearVariants[S ~[]E, E comparable]() []linearVariant[S, E] {
	return []linearVariant[S, E]{
		{name: "LinearSearch", f: search.LinearSearch[S]},
		{name: "LinearSearchManual", f: search.LinearSearchManual[S]},
	}
}

func TestLinearSearch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		in     []int
		target int
		want   int
	}{
		{
			name:   "first",
			in:     []int{3, 1, 4, 1, 5},
			target: 3,
			want:   0,
		},
		{
			name:   "middle",
			in:     []int{3, 1, 4, 1, 5},
			target: 4,
			want:   2,
		},
		{
			name:   "last",
			in:     []int{3, 1, 4, 1, 5},
			target: 5,
			want:   4,
		},
		{
			name:   "duplicate returns first",
			in:     []int{3, 1, 4, 1, 5},
			target: 1,
			want:   1,
		},
		{
			name:   "absent",
			in:     []int{3, 1, 4, 1, 5},
			target: 2,
			want:   -1,
		},
		{
			name:   "single",
			in:     []int{7},
			target: 7,
			want:   0,
		},
		{
			name:   "empty",
			in:     []int{},
			target: 1,
			want:   -1,
		},
		{
			name:   "nil",
			in:     nil,
			target: 1,
			want:   -1,
		},
	}

	for _, v := range linearVariants[[]int]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					if got := v.f(tt.in, tt.target); got != tt.want {
						t.Errorf("%s(%v, %v) = %v, want %v", v.name, tt.in, tt.target, got, tt.want)
					}
				})
			}
		})
	}
}

func TestLinearSearchSupportsComparableTypes(t *testing.T) {
	t.Parallel()
	type scores []int
	type point struct{ x, y int }
	testLinearSearch(t, scores{1, 3, 2}, 3, 1)
	testLinearSearch(t, []string{"a", "c", "b"}, "b", 2)
	testLinearSearch(t, []point{{x: 1, y: 2}, {x: 3, y: 4}}, point{x: 3, y: 4}, 1)
}

func testLinearSearch[S ~[]E, E comparable](t *testing.T, in S, target E, want int) {
	t.Helper()
	t.Run(reflect.TypeFor[S]().String(), func(t *testing.T) {
		t.Parallel()
		for _, v := range linearVariants[S]() {
			t.Run(v.name, func(t *testing.T) {
				t.Parallel()
				if got := v.f(in, target); got != want {
					t.Errorf("%s(%v, %v) = %v, want %v", v.name, in, target, got, want)
				}
			})
		}
	})
}

func TestLinearSearchUsesFloatEquality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		in     []float64
		target float64
		want   int
	}{
		{
			name:   "NaN never matches",
			in:     []float64{1, math.NaN(), 2},
			target: math.NaN(),
			want:   -1,
		},
		{
			name:   "negative zero matches zero",
			in:     []float64{math.Copysign(0, -1)},
			target: 0,
			want:   0,
		},
	}

	for _, v := range linearVariants[[]float64]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					if got := v.f(tt.in, tt.target); got != tt.want {
						t.Errorf("%s(%v, %v) = %v, want %v", v.name, tt.in, tt.target, got, tt.want)
					}
				})
			}
		})
	}
}
