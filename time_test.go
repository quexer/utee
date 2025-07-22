package utee

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeTracker", func() {
	var tt *TimeTracker

	BeforeEach(func() {
		tt = NewTimeTracker()
	})

	It("should track total time and steps correctly", func() {
		time.Sleep(10 * time.Millisecond)
		tt.Tick("step1")
		time.Sleep(20 * time.Millisecond)
		tt.Tick("step2")
		time.Sleep(10 * time.Millisecond)
		tt.Tick("step3")

		total := tt.Total()
		Expect(total).Should(BeNumerically(">", 0))

		s := tt.ToString()
		Expect(s).To(ContainSubstring("|"))
		Expect(s).To(ContainSubstring("-step1-"))
		Expect(s).To(ContainSubstring("step2"))
		Expect(s).To(ContainSubstring("step3"))
	})

	It("should contain all step names in ToString", func() {
		tt.Tick("first")
		tt.Tick("second")
		tt.Tick("third")

		steps := tt.ToString()
		Expect(steps).To(Equal("|-0s-first-0s-second-0s-third"))
	})
})
