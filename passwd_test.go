package utee_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"

	"github.com/quexer/utee"
)

var _ = Describe("Passwd, compatible with PHP Yii framework", func() {
	DescribeTable("VerifyPasswd",
		func(hash string) {
			ok := utee.VerifyPasswd("000000", hash)
			Ω(ok).To(BeTrue())
		},
		Entry("a,const=13", "$2y$13$PbOAJYUHxOIehQHNUGyE3uLqBXF4j1NKO4.iGDDYcsCQGr1pe4.Au"),
		Entry("b,const=13", "$2y$13$UaeOnuc9fANF7C3sBGI3buemqDuxjn8rOt.5oC4mvvtvL5gmyi76S"),
		Entry("c,const=13", "$2a$13$UcIbCAXU8i20hHmyS4ZiheISi7EZQy/Y.tpkISfMRcrSaeAkGjW5."),
		Entry("d,const=10", "$2a$10$Kn0DqGVtCBKPoY28IKBxi.C9hCjpd5uGBREXc5eR536myuBOAdhgq"),
	)
	It("PasswdHash", func() {
		hash, err := utee.PasswdHash("000000")
		Ω(err).To(Succeed())

		ok := utee.VerifyPasswd("000000", hash)
		Ω(ok).To(BeTrue())
	})
	It("get cost from hash", func() {
		n, err := bcrypt.Cost([]byte("$2a$10$Kn0DqGVtCBKPoY28IKBxi.C9hCjpd5uGBREXc5eR536myuBOAdhgq"))
		Ω(err).To(Succeed())
		Ω(n).To(Equal(10))
	})
})
