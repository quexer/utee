package utee

import (
	"github.com/samber/lo"
)

// SplitIntSlice split int slice into chunks
// Deprecated use lo.Chunk instead
func SplitIntSlice(src []int, chunkSize int) [][]int {
	return lo.Chunk(src, chunkSize)
}

// Min returns the smallest parameter
func Min(n ...int) int {
	return lo.Min(n)
}

// Max returns the biggest parameter
func Max(n ...int) int {
	return lo.Max(n)
}
