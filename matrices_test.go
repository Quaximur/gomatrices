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
	sizes := []int{10, 50, 100, 150, 190, 200, 500} // , 1000, 5000, 10000
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
}
