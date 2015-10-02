package rbo

import (
	"math"
	"testing"
)

func TestRBO(t *testing.T) {

	// from the reference implementation
	var tests = []struct {
		r1  []int
		r2  []int
		p   float64
		rbo float64
	}{
		{
			[]int{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'},
			[]int{'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'},
			0.95,
			0.771924,
		},
		{
			[]int{'g', 'a', 'f', 'c', 'z'},
			[]int{'a', 'b', 'c', 'd'},
			0.8,
			0.3786667,
		},
		{
			[]int{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'},
			[]int{'b', 'a', 'd', 'f', 'c', 'h'},
			0.8,
			0.6869504,
		},
		{
			[]int{'a', 'b', 'c', 'd', 'e'},
			[]int{'b', 'a', 'g', 'h', 'e', 'k', 'l', 'c'},
			0.9,
			0.6338971,
		},
	}

	for _, tt := range tests {
		if got := Calculate(tt.r1, tt.r2, tt.p); math.Abs(got-tt.rbo) > 0.000001 {
			t.Errorf("Calculate(%v,%v,%v)=%v, want %v\n", tt.r1, tt.r2, tt.p, got, tt.rbo)
		}
	}

}

func BenchmarkRBOState(b *testing.B) {

	l1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	l2 := []int{7, 5, 1, 8, 12, 4, 21, 27, 9, 26, 25, 24, 16, 0}

	for i := 0; i < b.N; i++ {
		Calculate(l1, l2, 0.98)
	}
}
