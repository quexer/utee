package utee

import (
	"regexp"
	"strings"
)

var (
	regHash = regexp.MustCompile(`[0-9a-zA-Z]{17,}`)
	regNum  = regexp.MustCompile(`\d{6,}`)
)

func CleanURI(s string) string {
	s = strings.Split(s, "?")[0]
	s = regHash.ReplaceAllString(s, ":hash")
	return regNum.ReplaceAllString(s, ":num")
}
