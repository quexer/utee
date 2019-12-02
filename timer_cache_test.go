package utee

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimerCache", func() {
	var tc *TimerCache
	BeforeEach(func() {
		tc = NewTimerCache(1, nil)
	})
	It("auto clear after expire", func() {
		By("brand new time cache should be empty", func() {
			Expect(tc.Len()).Should(Equal(0))
		})
		By("length should be changed after put 1", func() {
			tc.Put(1, "a")
			Expect(tc.Len()).Should(Equal(1))
			Expect(tc.Get(1)).Should(Equal("a"))
		})
		By("restore empty after expire", func() {
			Eventually(func() int { return tc.Len() }, 2).Should(Equal(0))
		})
	})
	It("keys", func() {
		Expect(len(tc.Keys())).Should(Equal(0))
		tc.Put(1, "a")
		Expect(tc.Keys()).Should(Equal([]interface{}{1}))
	})
	It("remove", func() {
		tc.Put(1, "a")
		tc.Remove(1)
		Expect(tc.Len()).Should(Equal(0))
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
})
