package utee

import (
	"github.com/samber/lo"
)

// SplitStringSlice split string slice into chunks
func SplitStringSlice(src []string, chunkSize int) [][]string {
	var out [][]string
	for {
		if len(src) == 0 {
			break
		}
		if len(src) < chunkSize {
			chunkSize = len(src)
		}
		out = append(out, src[0:chunkSize])
		src = src[chunkSize:]
	}
	return out
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
