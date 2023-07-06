package utee_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Slice", func() {
	Context("SplitSliceIntoN", func() {
		It("nil ", func() {
			a := utee.SplitSliceIntoN([]int(nil), 5)
			Ω(a).To(HaveLen(1))
			Ω(a[0]).To(BeNil())
		})
		It("one element", func() {
			a := utee.SplitSliceIntoN([]string{"a"}, 4)
			Ω(a).To(HaveLen(1))
			Ω(a[0][0]).To(Equal("a"))
		})
		DescribeTable("two element",
			func(maxSplit int) {
				a := utee.SplitSliceIntoN([]string{"a", "b"}, maxSplit)
				Ω(a).To(HaveLen(1))
				Ω(a[0]).To(HaveLen(2))
			},
			Entry("5", 5),
			Entry("1", 1),
		)
		It("4 => 3", func() {
			a := utee.SplitSliceIntoN([]string{"a", "b", "c", "d"}, 3)
			Ω(a).To(HaveLen(3))
			Ω(a[0]).To(HaveLen(2))
			Ω(a[1]).To(HaveLen(1))
			Ω(a[2]).To(HaveLen(1))
		})

		It("6 => 3", func() {
			a := utee.SplitSliceIntoN([]string{"a", "b", "c", "d", "e", "f"}, 3)
			Ω(a).To(HaveLen(3))
			Ω(a[2][1]).To(Equal("f"))
		})

		It("3 => 2", func() {
			a := utee.SplitSliceIntoN([]string{"a", "b", "c"}, 2)
			Ω(a).To(HaveLen(2))
			Ω(a[0]).To(Equal([]string{"a", "c"}))
			Ω(a[1]).To(Equal([]string{"b"}))
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
