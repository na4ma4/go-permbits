package permbits_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPermbits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Permbits Suite")
}
