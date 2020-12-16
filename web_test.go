package utee_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/quexer/utee"
)

var _ = Describe("Web", func() {
	Context("J", func() {
		var j utee.J
		BeforeEach(func() {
			j = utee.J{"name": "a", "id": 5}
		})
		It("ToString", func() {
			Ω(j.ToString()).Should(MatchJSON(`{"id":5, "name":"a"}`))
		})
		It("ToReader", func() {
			b, _ := ioutil.ReadAll(j.ToReader())
			Ω(string(b)).Should(MatchJSON(`{"id":5, "name":"a"}`))
		})
	})
})
