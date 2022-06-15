package utee_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("String", func() {
	DescribeTable("SplitStringSlice",
		func(cap, n int) {
			l := utee.SplitStringSlice([]string{"a", "b", "c", "d", "e", "f"}, cap)
			Ω(l).To(HaveLen(n))
			Ω(len(l[0])).To(BeNumerically(">", 0))
		},
		Entry("4-2", 4, 2),
		Entry("2-3", 2, 3),
		Entry("6-1", 6, 1),
		Entry("100-1", 6, 1),
	)
	Context("SplitStringSliceIntoN", func() {
		It("nil ", func() {
			a := utee.SplitStringSliceIntoN(nil, 5)
			Ω(a).To(HaveLen(1))
			Ω(a[0]).To(BeNil())
		})
		It("one element", func() {
			a := utee.SplitStringSliceIntoN([]string{"a"}, 4)
			Ω(a).To(HaveLen(1))
			Ω(a[0][0]).To(Equal("a"))
		})
		DescribeTable("two element",
			func(maxSplit int) {
				a := utee.SplitStringSliceIntoN([]string{"a", "b"}, maxSplit)
				Ω(a).To(HaveLen(1))
				Ω(a[0]).To(HaveLen(2))
			},
			Entry("5", 5),
			Entry("1", 1),
		)
		It("4 => 3", func() {
			a := utee.SplitStringSliceIntoN([]string{"a", "b", "c", "d"}, 3)
			Ω(a).To(HaveLen(3))
			Ω(a[0]).To(HaveLen(2))
			Ω(a[1]).To(HaveLen(1))
			Ω(a[2]).To(HaveLen(1))
		})

		It("6 => 3", func() {
			a := utee.SplitStringSliceIntoN([]string{"a", "b", "c", "d", "e", "f"}, 3)
			Ω(a).To(HaveLen(3))
			Ω(a[2][1]).To(Equal("f"))
		})

		It("3 => 2", func() {
			a := utee.SplitStringSliceIntoN([]string{"a", "b", "c"}, 2)
			Ω(a).To(HaveLen(2))
			Ω(a[0]).To(Equal([]string{"a", "c"}))
			Ω(a[1]).To(Equal([]string{"b"}))
		})

	})
})
