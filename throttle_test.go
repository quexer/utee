package utee

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Throttle", func() {
	It("Number", func() {
		latch := NewThrottle(100)

		latch.Acquire()
		latch.Acquire()

		Ω(latch.Current()).To(Equal(2))
		Ω(latch.Available()).To(Equal(98))
		latch.Release()
		Ω(latch.Current()).To(Equal(1))
		Ω(latch.Available()).To(Equal(99))
	})
	It("Wait", func() {
		latch := NewThrottle(3)
		go func() {
			time.Sleep(2 * time.Second)
			latch.Release()
			latch.Release()
		}()
		latch.Acquire()
		latch.Acquire()
		latch.Wait()
	})
})
