package gomatrices

import (
	"fmt"
	"slices"
	"testing"
)

type multiplyFunc func(m1, m2 [][]int) ([][]int, error)

var multipliers = map[string]multiplyFunc{
	"Small":  MultiplySmall[int],
	"Medium": MultiplyMedium[int],
	"Huge":   MultiplyHuge[int],
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		m1, m2   [][]int
		expected [][]int
	}{
		{
			name:     "3x3 * 3x3",
			m1:       generateMatrix(3), // 3x3 with values 1..9
			m2:       generateMatrix(3),
			expected: [][]int{{30, 36, 42}, {66, 81, 96}, {102, 126, 150}},
		},
		{
			name:     "2x3 * 3x4",
			m1:       [][]int{{1, 2, 3}, {4, 5, 6}},
			m2:       [][]int{{1, 0, 0, 1}, {0, 1, 0, 1}, {0, 0, 1, 1}},
			expected: [][]int{{1, 2, 3, 6}, {4, 5, 6, 15}},
		},
		{
			name:     "3x2 * 2x1",
			m1:       [][]int{{1, 2}, {3, 4}, {5, 6}},
			m2:       [][]int{{7}, {8}},
			expected: [][]int{{23}, {53}, {83}},
		},
		{
			name:     "1x3 * 3x1",
			m1:       [][]int{{1, 2, 3}},
			m2:       [][]int{{4}, {5}, {6}},
			expected: [][]int{{32}},
		},
	}

	for funcName, multiply := range multipliers {
		for _, tt := range tests {
			t.Run(funcName+"/"+tt.name, func(t *testing.T) {
				resultMatrix, err := multiply(tt.m1, tt.m2)
				if err != nil {
					t.Fatalf("Unexpected error '%s'. Should be no error", err)
				}

				if !slices.EqualFunc(resultMatrix, tt.expected, slices.Equal) {
					t.Fatalf("Unexpected matrix result '%v'. Should be '%v'", resultMatrix, tt.expected)
				}
			})
		}
	}
}

func TestMultiplyInvalidShapes(t *testing.T) {
	m23 := [][]int{{1, 2, 3}, {4, 5, 6}}

	tests := []struct {
		name   string
		m1, m2 [][]int
	}{
		{"nil matrix1", nil, m23},
		{"empty matrix2", m23, [][]int{}},
		{"2x3 * 4x2 (too many rows)", m23, [][]int{{1, 0}, {0, 1}, {1, 1}, {9, 9}}},
		{"2x3 * 2x2 (too few rows)", m23, [][]int{{1, 0}, {0, 1}}},
	}

	for funcName, multiply := range multipliers {
		for _, tt := range tests {
			t.Run(funcName+"/"+tt.name, func(t *testing.T) {
				resultMatrix, err := multiply(tt.m1, tt.m2)
				if err == nil {
					t.Fatalf("Expected error, got result '%v'", resultMatrix)
				}
			})
		}
	}
}

func BenchmarkMultiply(b *testing.B) {
	sizes := []int{10, 50, 100, 150, 190, 200, 500, 1000, 5000, 10000}
	for _, size := range sizes {
		m1 := generateMatrix(size)
		m2 := generateMatrix(size)

		// Small is single-threaded O(n^3), 5000+ would take hours
		if size <= 1000 {
			b.Run(fmt.Sprintf("Small_%d", size), func(b *testing.B) {
				for b.Loop() {
					MultiplySmall(m1, m2)
				}
			})
		}

		b.Run(fmt.Sprintf("Medium_%d", size), func(b *testing.B) {
			for b.Loop() {
				MultiplyMedium(m1, m2)
			}
		})

		b.Run(fmt.Sprintf("Huge_%d", size), func(b *testing.B) {
			for b.Loop() {
				MultiplyHuge(m1, m2)
			}
		})
	}

	// goos: windows
	// goarch: amd64
	// pkg: github.com/quaximur/testo/optimath
	// cpu: Intel(R) Core(TM) i7-14700HX
	// BenchmarkMultiply/Small_10-28         	  768358	      1368 ns/op	    1040 B/op	      11 allocs/op
	// BenchmarkMultiply/Medium_10-28        	  168678	      6852 ns/op	    2336 B/op	      32 allocs/op
	// BenchmarkMultiply/Huge_10-28          	   52740	     22398 ns/op	    5890 B/op	      52 allocs/op
	// BenchmarkMultiply/Small_50-28         	    8637	    146577 ns/op	   22080 B/op	      51 allocs/op
	// BenchmarkMultiply/Medium_50-28        	   25341	     46971 ns/op	   28496 B/op	     152 allocs/op
	// BenchmarkMultiply/Huge_50-28          	   15534	     75536 ns/op	   48305 B/op	     132 allocs/op
	// BenchmarkMultiply/Small_100-28        	    1033	   1057842 ns/op	   92288 B/op	     101 allocs/op
	// BenchmarkMultiply/Medium_100-28       	    6596	    182706 ns/op	  105112 B/op	     302 allocs/op
	// BenchmarkMultiply/Huge_100-28         	    4880	    256904 ns/op	  189280 B/op	     232 allocs/op
	// BenchmarkMultiply/Small_150-28        	     302	   3842535 ns/op	  196096 B/op	     151 allocs/op
	// BenchmarkMultiply/Medium_150-28       	    2468	    484574 ns/op	  215334 B/op	     452 allocs/op
	// BenchmarkMultiply/Huge_150-28         	    1760	    581890 ns/op	  397282 B/op	     332 allocs/op
	// BenchmarkMultiply/Small_190-28        	     148	   8017443 ns/op	  296704 B/op	     191 allocs/op
	// BenchmarkMultiply/Medium_190-28       	    1416	    823064 ns/op	  321058 B/op	     572 allocs/op
	// BenchmarkMultiply/Huge_190-28         	    1438	    991454 ns/op	  598906 B/op	     412 allocs/op
	// BenchmarkMultiply/Small_200-28        	     123	   9953143 ns/op	  363264 B/op	     201 allocs/op
	// BenchmarkMultiply/Medium_200-28       	    1288	    930111 ns/op	  388881 B/op	     602 allocs/op
	// BenchmarkMultiply/Huge_200-28         	    1275	    945423 ns/op	  732033 B/op	     433 allocs/op
	// BenchmarkMultiply/Small_500-28        	       4	 260439975 ns/op	 2060288 B/op	     501 allocs/op
	// BenchmarkMultiply/Medium_500-28       	      74	  15306885 ns/op	 2124480 B/op	    1502 allocs/op
	// BenchmarkMultiply/Huge_500-28         	     121	   9934171 ns/op	 4129148 B/op	    1032 allocs/op
	// BenchmarkMultiply/Small_1000-28       	       1	3964509400 ns/op	 8216576 B/op	    1001 allocs/op
	// BenchmarkMultiply/Medium_1000-28      	       4	 270904850 ns/op	 8363696 B/op	    3029 allocs/op
	// BenchmarkMultiply/Huge_1000-28        	      16	  76295488 ns/op	16445049 B/op	    2032 allocs/op
	// PASS
	// ok  	github.com/quaximur/testo/optimath	43.734s
}
