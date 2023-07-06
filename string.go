package utee

import (
	"github.com/samber/lo"
)

// Truncate , truncate string as []rune
// make sure the rune count of result is not more than maxLen
func Truncate(s string, maxLen uint) string {
	return string(lo.Subset([]rune(s), 0, maxLen))
}
