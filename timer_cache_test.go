package utee

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimerCache", func() {
	var tc *TimerCache
	BeforeEach(func() {
		tc = NewTimerCache(1, nil)
	})
	It("auto clear after expire", func() {
		By("brand new time cache should be empty", func() {
			Ω(tc.Len()).To(Equal(0))
		})
		By("length should be changed after put 1", func() {
			tc.Put(1, "a")
			Ω(tc.Len()).To(Equal(1))
			Ω(tc.Get(1)).To(Equal("a"))
		})
		By("restore empty after expire", func() {
			Eventually(func() int { return tc.Len() }, 2).Should(Equal(0))
		})
	})
	It("keys", func() {
		Ω(len(tc.Keys())).To(Equal(0))
		tc.Put(1, "a")
		Ω(tc.Keys()).To(Equal([]interface{}{1}))
	})
	It("remove", func() {
		tc.Put(1, "a")
		tc.Remove(1)
		Ω(tc.Len()).To(Equal(0))
	})
	It("expire", func() {
		var key int
		var val string
		tc = NewTimerCache(1, func(k, v interface{}) {
			key = k.(int)
			val = v.(string)
		})
		tc.Put(1, "a")
		Eventually(func() int { return key }, 2).Should(Equal(1))
		Eventually(func() string { return val }, 2).Should(Equal("a"))
	})
	It("ttl", func() {
		tc = NewTimerCache(2, nil)
		tc.Put(1, "a")
		Ω(tc.TTL(1)).To(BeEquivalentTo(2))
		Eventually(func() int64 { return tc.TTL(1) }, 3).Should(BeEquivalentTo(0))
	})
})
