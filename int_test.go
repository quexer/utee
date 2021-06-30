package utee_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Int", func() {
	Context("SplitIntSlice", func() {
		It("no split", func() {
			l := utee.SplitIntSlice([]int{1, 2, 3, 4, 5}, 100)
			Ω(l).To(HaveLen(1))
			Ω(l[0]).To(Equal([]int{1, 2, 3, 4, 5}))
		})
		It("split", func() {
			l := utee.SplitIntSlice([]int{1, 2, 3, 4, 5}, 2)
			Ω(l).To(HaveLen(3))
			Ω(l).To(Equal([][]int{
				{1, 2}, {3, 4}, {5},
			}))
		})
	})
	DescribeTable("Min",
		func(src []int, min int) {
			n := utee.Min(src...)
			Ω(n).To(Equal(min))
		},
		Entry("a", []int{7, 5, 1}, 1),
		Entry("b", []int{1, 3, 7, 5}, 1),
	)
	DescribeTable("Max",
		func(src []int, max int) {
			n := utee.Max(src...)
			Ω(n).To(Equal(max))
		},
		Entry("a", []int{7, 5, 1}, 7),
		Entry("b", []int{1, 3, 7, 5}, 7),
	)
})
