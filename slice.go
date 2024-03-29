package utee

import (
	"cmp"

	"github.com/samber/lo"
)

// Min returns the smallest parameter
func Min[T cmp.Ordered](n ...T) T {
	return lo.Min(n)
}

// Max returns the biggest parameter
func Max[T cmp.Ordered](n ...T) T {
	return lo.Max(n)
}

// SplitSliceIntoN split slice into several parts, no more than n
func SplitSliceIntoN[T any](a []T, n int) [][]T {
	if len(a) < n || n == 1 {
		return [][]T{a}
	}

	result := make([][]T, n)

	for i, s := range a {
		idx := i % n
		result[idx] = append(result[idx], s)
	}

	return result
}
