package utee_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Tick", func() {
	It("TickToTime", func() {
		now := time.Now()
		t := utee.NewTick(now).ToTime()
		Ω(t).To(BeTemporally("~", now, time.Millisecond))
	})
})
