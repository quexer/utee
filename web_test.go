package utee_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Web", func() {
	It("J.ToString()", func() {
		s := utee.J{"name": "a", "id": 5}.ToString()
		Ω(s).Should(MatchJSON(`{"id":5, "name":"a"}`))
	})
})
