package utee

import (
	"github.com/samber/lo"
)

// SplitStringSlice split string slice into chunks
// Deprecated use lo.Chunk instead
func SplitStringSlice(src []string, chunkSize int) [][]string {
	return lo.Chunk(src, chunkSize)
}

// SplitStringSliceIntoN split a into several parts, no more than n
func SplitStringSliceIntoN(a []string, n int) [][]string {
	if len(a) < n || n == 1 {
		return [][]string{a}
	}

	result := make([][]string, n)
	for i, s := range a {
		idx := i % n
		result[idx] = append(result[idx], s)
	}
	return result
}

// Truncate , truncate string as []rune
// make sure the rune count of result is not more than maxLen
func Truncate(s string, maxLen uint) string {
	return string(lo.Subset([]rune(s), 0, maxLen))
}
