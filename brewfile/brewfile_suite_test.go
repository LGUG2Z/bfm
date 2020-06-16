package brewfile_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBrewfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brewfile Suite")
}

var testPath = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/LGUG2Z/bfm/testData")

var _ = BeforeSuite(func() {
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		Expect(os.Mkdir(testPath, os.ModePerm)).To(Succeed())
	}
})

var _ = AfterSuite(func() {
	if _, err := os.Stat(testPath); err == nil {
		Expect(os.Remove(testPath)).To(Succeed())
	}
})
