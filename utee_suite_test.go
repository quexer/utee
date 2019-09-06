package utee_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtee(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utee Suite")
}
