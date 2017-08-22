package brewfile_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBrewfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brewfile Suite")
}
