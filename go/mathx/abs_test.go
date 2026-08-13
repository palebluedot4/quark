package mathx_test

import (
	"math"
	"math/rand/v2"
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"

	"github.com/palebluedot4/quark/go/mathx"
)

type variant[T any] struct {
	name string
	f    func(T) T
}

func signedVariants[T constraints.Signed]() []variant[T] {
	return []variant[T]{
		{name: "AbsSigned", f: mathx.AbsSigned[T]},
		{name: "AbsSignedBitwise", f: mathx.AbsSignedBitwise[T]},
	}
}

func floatVariants[T constraints.Float]() []variant[T] {
	return []variant[T]{
		{name: "AbsFloat", f: mathx.AbsFloat[T]},
		{name: "AbsFloatSignbit", f: mathx.AbsFloatSignbit[T]},
		{name: "AbsFloatBitwise", f: mathx.AbsFloatBitwise[T]},
	}
}

func TestAbsSigned(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   int
		want int
	}{
		{
			name: "positive",
			in:   5,
			want: 5,
		},
		{
			name: "negative",
			in:   -5,
			want: 5,
		},
		{
			name: "zero",
			in:   0,
			want: 0,
		},
		{
			name: "maximum",
			in:   math.MaxInt,
			want: math.MaxInt,
		},
		{
			name: "minimum wraps",
			in:   math.MinInt,
			want: math.MinInt,
		},
	}

	for _, v := range signedVariants[int]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					if got := v.f(tt.in); got != tt.want {
						t.Errorf("%s(%v) = %v, want %v", v.name, tt.in, got, tt.want)
					}
				})
			}
		})
	}
}

func TestAbsSignedAtEveryWidth(t *testing.T) {
	t.Parallel()
	testSignedWidth[int8](t, math.MinInt8)
	testSignedWidth[int16](t, math.MinInt16)
	testSignedWidth[int32](t, math.MinInt32)
	testSignedWidth[int64](t, math.MinInt64)
}

func testSignedWidth[T constraints.Signed](t *testing.T, minimum T) {
	t.Helper()
	tests := []struct {
		name string
		in   T
		want T
	}{
		{
			name: "negative",
			in:   -5,
			want: 5,
		},
		{
			name: "minimum wraps",
			in:   minimum,
			want: minimum,
		},
	}

	t.Run(reflect.TypeFor[T]().Name(), func(t *testing.T) {
		t.Parallel()
		for _, v := range signedVariants[T]() {
			t.Run(v.name, func(t *testing.T) {
				t.Parallel()
				for _, tt := range tests {
					t.Run(tt.name, func(t *testing.T) {
						t.Parallel()
						if got := v.f(tt.in); got != tt.want {
							t.Errorf("%s(%v) = %v, want %v", v.name, tt.in, got, tt.want)
						}
					})
				}
			})
		}
	})
}

func TestAbsSignedSupportsNamedTypes(t *testing.T) {
	t.Parallel()
	type temperature int16
	for _, v := range signedVariants[temperature]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			got := v.f(-5)
			if want := temperature(5); got != want {
				t.Errorf("%s(temperature(-5)) = %v, want %v", v.name, got, want)
			}
		})
	}
}

func TestAbsFloat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   float64
		want float64
	}{
		{
			name: "positive",
			in:   5.5,
			want: 5.5,
		},
		{
			name: "negative",
			in:   -5.5,
			want: 5.5,
		},
		{
			name: "zero",
			in:   0,
			want: 0,
		},
		{
			name: "positive infinity",
			in:   math.Inf(1),
			want: math.Inf(1),
		},
		{
			name: "negative infinity",
			in:   math.Inf(-1),
			want: math.Inf(1),
		},
		{
			name: "largest",
			in:   -math.MaxFloat64,
			want: math.MaxFloat64,
		},
		{
			name: "smallest subnormal",
			in:   -math.SmallestNonzeroFloat64,
			want: math.SmallestNonzeroFloat64,
		},
	}

	for _, v := range floatVariants[float64]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					if got := v.f(tt.in); got != tt.want {
						t.Errorf("%s(%v) = %v, want %v", v.name, tt.in, got, tt.want)
					}
				})
			}
		})
	}
}

func TestAbsFloatNormalizesSignAtEveryWidth(t *testing.T) {
	t.Parallel()
	testFloatSign[float32](t)
	testFloatSign[float64](t)
}

func testFloatSign[T constraints.Float](t *testing.T) {
	t.Helper()
	tests := []struct {
		name    string
		in      T
		wantNaN bool
	}{
		{
			name:    "positive zero",
			in:      0,
			wantNaN: false,
		},
		{
			name:    "negative zero",
			in:      T(math.Copysign(0, -1)),
			wantNaN: false,
		},
		{
			name:    "positive NaN",
			in:      T(math.NaN()),
			wantNaN: true,
		},
		{
			name:    "negative NaN",
			in:      T(math.Copysign(math.NaN(), -1)),
			wantNaN: true,
		},
	}

	t.Run(reflect.TypeFor[T]().Name(), func(t *testing.T) {
		t.Parallel()
		for _, v := range floatVariants[T]() {
			t.Run(v.name, func(t *testing.T) {
				t.Parallel()
				for _, tt := range tests {
					t.Run(tt.name, func(t *testing.T) {
						t.Parallel()
						got := v.f(tt.in)
						if math.Signbit(float64(got)) {
							t.Errorf("%s(%v) = %v, want a positive result", v.name, tt.in, got)
						}
						if tt.wantNaN && !math.IsNaN(float64(got)) {
							t.Errorf("%s(%v) = %v, want NaN", v.name, tt.in, got)
						}
					})
				}
			})
		}
	})
}

func TestAbsFloatPreservesFloat32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   float32
		want float32
	}{
		{
			name: "not representable in binary",
			in:   -1.1,
			want: 1.1,
		},
		{
			name: "largest",
			in:   -math.MaxFloat32,
			want: math.MaxFloat32,
		},
		{
			name: "smallest subnormal",
			in:   -math.SmallestNonzeroFloat32,
			want: math.SmallestNonzeroFloat32,
		},
	}

	for _, v := range floatVariants[float32]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					if got := v.f(tt.in); got != tt.want {
						t.Errorf("%s(float32(%v)) = %v, want %v", v.name, tt.in, got, tt.want)
					}
				})
			}
		})
	}
}

func TestAbsFloatSupportsNamedTypes(t *testing.T) {
	t.Parallel()
	type measurement float32
	for _, v := range floatVariants[measurement]() {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			got := v.f(-5.5)
			if want := measurement(5.5); got != want {
				t.Errorf("%s(measurement(-5.5)) = %v, want %v", v.name, got, want)
			}
		})
	}
}

const inputSize = 1 << 10

func BenchmarkAbsSigned(b *testing.B) {
	r := rand.New(rand.NewPCG(42, 0))
	in := make([]int, inputSize)
	for i := range in {
		in[i] = r.IntN(2001) - 1000
	}
	benchmarkVariants(b, signedVariants[int](), in)
}

func BenchmarkAbsFloat64(b *testing.B) {
	r := rand.New(rand.NewPCG(42, 0))
	in := make([]float64, inputSize)
	for i := range in {
		in[i] = r.Float64()*2000 - 1000
	}
	benchmarkVariants(b, floatVariants[float64](), in)
}

func BenchmarkAbsFloat32(b *testing.B) {
	r := rand.New(rand.NewPCG(42, 0))
	in := make([]float32, inputSize)
	for i := range in {
		in[i] = r.Float32()*2000 - 1000
	}
	benchmarkVariants(b, floatVariants[float32](), in)
}

func benchmarkVariants[T any](b *testing.B, variants []variant[T], in []T) {
	b.Helper()
	for _, v := range variants {
		b.Run(v.name, func(b *testing.B) {
			i := 0
			for b.Loop() {
				v.f(in[i&(len(in)-1)])
				i++
			}
		})
	}
}
