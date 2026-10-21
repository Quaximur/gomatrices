package gomatrices

import (
	"errors"
	"runtime"
	"sync"

	"golang.org/x/exp/constraints"
)

type numeric interface {
	constraints.Integer | constraints.Float
}

// Checks that [N, M] * [M, K] is possible. Only first rows are inspected, ragged matrices are not detected
func checkShapes[T numeric](matrix1, matrix2 [][]T) error {
	if len(matrix1) == 0 {
		return errors.New("matrix1 is nil or empty")
	}

	if len(matrix2) == 0 {
		return errors.New("matrix2 is nil or empty")
	}

	if len(matrix1[0]) != len(matrix2) {
		return errors.New("matrix1 columns len != matrix2 rows len")
	}

	return nil
}

// Multiply [N, M] * [M, K] matrices (len ~1-50). CAUTION: No ragged matrices guard implemented
func MultiplySmall[T numeric](m1, m2 [][]T) ([][]T, error) {
	if err := checkShapes(m1, m2); err != nil {
		return nil, err
	}

	rows, cols, inner := len(m1), len(m2[0]), len(m2)
	result := make([][]T, rows)
	for i := range result {
		result[i] = make([]T, cols)
		for j := range cols {
			var sum T
			for k := range inner {
				sum += m1[i][k] * m2[k][j]
			}
			result[i][j] = sum
		}
	}
	return result, nil
}

// Multiply [N, M] * [M, K] matrices (len ~50-190). CAUTION: No ragged matrices guard implemented
func MultiplyMedium[T numeric](matrix1, matrix2 [][]T) ([][]T, error) {
	if err := checkShapes(matrix1, matrix2); err != nil {
		return nil, err
	}

	len1 := len(matrix1)
	resMat := make([][]T, len1)
	resultCols := len(matrix2[0])

	var wg sync.WaitGroup
	wg.Add(len1)

	for i, r1 := range matrix1 {
		resMat[i] = make([]T, resultCols)

		go func(row1 []T, mat2 [][]T, resColumn []T) {
			defer wg.Done()
			for j2 := range resultCols {
				var i1j2Res T
				for i1, re1 := range row1 {
					i1j2Res += re1 * mat2[i1][j2]
				}
				resColumn[j2] = i1j2Res
			}
		}(r1, matrix2, resMat[i])
	}

	wg.Wait()

	return resMat, nil
}

// Multiply [N, M] * [M, K] matrices (len 190+). CAUTION: No ragged matrices guard implemented
func MultiplyHuge[T numeric](matrix1, matrix2 [][]T) ([][]T, error) {
	if err := checkShapes(matrix1, matrix2); err != nil {
		return nil, err
	}

	rows1 := len(matrix1)
	rows2 := len(matrix2)
	cols1 := len(matrix1[0])
	cols2 := len(matrix2[0])

	m2T := make([][]T, cols2)
	for j := range cols2 {
		m2T[j] = make([]T, rows2)
		for i := range rows2 {
			m2T[j][i] = matrix2[i][j]
		}
	}

	resMat := make([][]T, rows1)
	for i := range resMat {
		resMat[i] = make([]T, cols2)
	}

	// GOMAXPROCS respects -cpu flag and container CPU limits, unlike NumCPU
	numWorkers := runtime.GOMAXPROCS(0)
	jobs := make(chan int, rows1)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for range numWorkers {
		go func() {
			defer wg.Done()
			for i := range jobs {
				row1 := matrix1[i]
				for j := range cols2 {
					col2 := m2T[j]
					var sum T
					for k := range cols1 {
						sum += row1[k] * col2[k]
					}
					resMat[i][j] = sum
				}
			}
		}()
	}

	for i := range rows1 {
		jobs <- i
	}
	
	close(jobs)
	wg.Wait()

	return resMat, nil
}

func generateMatrix(size int) [][]int {
	matrix := make([][]int, size)
	for i := range size {
		matrix[i] = make([]int, size)
		for j := range size {
			matrix[i][j] = i*size + j + 1
		}
	}
	return matrix
}
