package utee

import (
	"regexp"
	"strings"
)

var (
	regHash = regexp.MustCompile(`[0-9a-zA-Z]{17,}`)
	regNum  = regexp.MustCompile(`\d{6,}`)
	regCrop = regexp.MustCompile(`/\d{2,}x\d{2,}$`)
)

func CleanURI(s string) string {
	s = strings.Split(s, "?")[0]
	s = regHash.ReplaceAllString(s, ":hash")
	s = regCrop.ReplaceAllString(s, "/:crop")
	return regNum.ReplaceAllString(s, ":num")
}
