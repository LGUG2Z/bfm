package brew_test

import (
	. "github.com/lgug2z/bfm/brew"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("InfoCache", func() {
	Describe("With a fresh unpopulated instance", func() {
		It("Should read the cached info from a file on disk", func() {
			expected := InfoCache([]Info{Info{Name: "a2ps", FullName: "a2ps", Desc: "Any-to-PostScript filter"}})

			actual := InfoCache([]Info{})
			error := actual.Read(fmt.Sprintf("%s/src/github.com/lgug2z/bfm/testData/test.json", os.Getenv("GOPATH")))
			Expect(error).To(BeNil())
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("With a populated InfoCache", func() {
		It("Should find and return the Info of a package", func() {
			cache := InfoCache([]Info{Info{FullName: "a"}})
			expected := Info{FullName: "a"}
			actual, error := cache.Find("a")

			Expect(error).To(BeNil())
			Expect(actual).To(Equal(expected))
		})

		It("Should return an error if a package cannot be found", func() {
			cache := InfoCache([]Info{Info{FullName: "a"}})
			_, error := cache.Find("b")

			Expect(error).ToNot(BeNil())
		})
	})
})
