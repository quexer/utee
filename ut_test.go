package utee_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Ut", func() {
	It("TickToTime", func() {
		now := time.Now()
		t := utee.NewTick(now).ToTime()
		Ω(t).To(BeTemporally("~", now, time.Millisecond))
	})

	It("Truncate: no truncate", func() {
		Ω(utee.Truncate("中文test", 10)).To(Equal("中文test"))
	})

})
