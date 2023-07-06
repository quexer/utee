package utee_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("String", func() {
	DescribeTable("Truncate",
		func(n int, result string) {
			s := "ab中文😜好"
			Ω(utee.Truncate(s, uint(n))).To(Equal(result))
		},
		Entry(nil, 0, ""),
		Entry(nil, 1, "a"),
		Entry(nil, 4, "ab中文"),
		Entry(nil, 5, "ab中文😜"),
		Entry(nil, 6, "ab中文😜好"),
		Entry(nil, 100, "ab中文😜好"),
	)
})
