package brewfile_test

import (
	. "github.com/lgug2z/bfm/brewfile"

	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packages", func() {
	Describe("When given a path to a Brewfile", func() {
		var (
			packages           Packages
			bf, info, contents string
		)

		BeforeEach(func() {
			bf = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")
			contents = `
tap 'homebrew/bundle'
brew 'a2ps'
tap 'homebrew/core'
cask 'google-chrome'
mas 'Xcode', id: 497799835
cask 'firefox'
# some comment
`
			ioutil.WriteFile(bf, []byte(contents), 0644)
			info = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/test.json")
		})

		AfterEach(func() {
			os.Remove(bf)
		})

		It("Reads, separates and stores packages from the Brewfile", func() {
			bf := fmt.Sprintf(bf)
			packages.FromBrewfile(bf)
			expectedPackages := Packages{
				Tap:  []string{"tap 'homebrew/bundle'", "tap 'homebrew/core'"},
				Brew: []string{"brew 'a2ps'"},
				Cask: []string{"cask 'firefox'", "cask 'google-chrome'"},
				Mas:  []string{"mas 'Xcode', id: 497799835"},
			}

			Expect(packages).To(Equal(expectedPackages))
		})
	})

	Describe("When popualted with packages", func() {
		It("Produces a byte representation of the contents to be written to disk", func() {

			packages := Packages{
				Tap:  []string{"tap 'homebrew/bundle'", "tap 'homebrew/core'"},
				Brew: []string{"brew 'a2ps'"},
				Cask: []string{"cask 'firefox'", "cask 'google-chrome'"},
				Mas:  []string{"mas 'Xcode', id: 497799835"},
			}

			actual := packages.Bytes()
			expected := []byte(`tap 'homebrew/bundle'
tap 'homebrew/core'

brew 'a2ps'

cask 'firefox'
cask 'google-chrome'

mas 'Xcode', id: 497799835`)
			Expect(actual).To(Equal(expected))
		})
	})

})
