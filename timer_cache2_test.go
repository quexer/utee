package utee

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimerCache2", func() {
	var tc *TimerCache2[int, string]
	BeforeEach(func() {
		tc = NewTimerCache2[int, string](1, nil)
	})
	It("auto clear after expire", func() {
		By("brand new time cache should be empty", func() {
			Ω(tc.Len()).To(Equal(0))
		})
		By("length should be changed after put 1", func() {
			tc.Put(1, "a")
			Ω(tc.Len()).To(Equal(1))
			s, b := tc.Get(1)
			Ω(b).To(BeTrue())
			Ω(s).To(Equal("a"))
		})
		By("restore empty after expire", func() {
			Eventually(func() int { return tc.Len() }, 2).Should(Equal(0))
		})
	})
	It("keys", func() {
		Ω(len(tc.Keys())).To(Equal(0))
		tc.Put(1, "a")
		Ω(tc.Keys()).To(Equal([]int{1}))
	})
	It("remove", func() {
		tc.Put(1, "a")
		tc.Remove(1)
		Ω(tc.Len()).To(Equal(0))
	})
	It("expire", func() {
		var key int
		var val string
		tc = NewTimerCache2[int, string](1, func(k int, v string) {
			key = k
			val = v
		})
		tc.Put(1, "a")
		Eventually(func() int { return key }, 2).Should(Equal(1))
		Eventually(func() string { return val }, 2).Should(Equal("a"))
	})
	It("ttl", func() {
		tc = NewTimerCache2[int, string](2, nil)
		tc.Put(1, "a")
		ttl, b := tc.TTL(1)
		Ω(b).To(BeTrue())
		Ω(ttl).To(BeEquivalentTo(2))
		Eventually(func() int64 {
			ttl, _ := tc.TTL(1)
			return ttl
		}, 3).Should(BeZero())
	})
})
