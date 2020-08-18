package utee_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Ut", func() {
	It("TickToTime", func() {
		now := time.Now()
		t := utee.TickToTime(utee.Tick(now))
		Expect(t).To(BeTemporally("~", now, time.Millisecond))
	})

	It("Truncate: no truncate", func() {
		Expect(utee.Truncate("中文test", 10)).To(Equal("中文test"))
	})

	Context("SplitSlice", func() {
		It("nil ", func() {
			a := utee.SplitSlice(nil, 5)
			Expect(a).To(HaveLen(1))
			Expect(a[0]).To(BeNil())
		})
		It("one element", func() {
			a := utee.SplitSlice([]string{"a"}, 4)
			Expect(a).To(HaveLen(1))
			Expect(a[0][0]).To(Equal("a"))
		})
		DescribeTable("two element",
			func(maxSplit int) {
				a := utee.SplitSlice([]string{"a", "b"}, maxSplit)
				Expect(a).To(HaveLen(1))
				Expect(a[0]).To(HaveLen(2))
			},
			Entry("5", 5),
			Entry("1", 1),
		)
		It("4 => 3", func() {
			a := utee.SplitSlice([]string{"a", "b", "c", "d"}, 3)
			Expect(a).To(HaveLen(3))
			Expect(a[0]).To(HaveLen(2))
			Expect(a[1]).To(HaveLen(1))
			Expect(a[2]).To(HaveLen(1))
		})

		It("6 => 3", func() {
			a := utee.SplitSlice([]string{"a", "b", "c", "d", "e", "f"}, 3)
			Expect(a).To(HaveLen(3))
			Expect(a[2][1]).To(Equal("f"))
		})

		It("3 => 2", func() {
			a := utee.SplitSlice([]string{"a", "b", "c"}, 2)
			Expect(a).To(HaveLen(2))
			Expect(a[0]).To(Equal([]string{"a", "c"}))
			Expect(a[1]).To(Equal([]string{"b"}))
		})

	})
})
