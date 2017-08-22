package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	"fmt"
	"io/ioutil"
	"os"

	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clean", func() {
	Describe("When the command is called", func() {

		var (
			cache              brew.InfoCache
			packages           brewfile.Packages
			bf, info, contents string
			flags              Flags
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

		It("Should read in the packages currently in the Brewfile", func() {
			expectedPackages := brewfile.Packages{
				Tap:  []string{"tap 'homebrew/bundle'", "tap 'homebrew/core'"},
				Brew: []string{"brew 'a2ps'"},
				Cask: []string{"cask 'firefox'", "cask 'google-chrome'"},
				Mas:  []string{"mas 'Xcode', id: 497799835"},
			}

			flags = Flags{false}

			Clean([]string{}, &packages, cache, bf, info, flags)

			Expect(packages).To(Equal(expectedPackages))
		})

		It("Should write out a new Brewfile in alphabetical order split into tap, brew, cask and mas sections", func() {
			expectedContents := `tap 'homebrew/bundle'
tap 'homebrew/core'

brew 'a2ps'

cask 'firefox'
cask 'google-chrome'

mas 'Xcode', id: 497799835`

			flags = Flags{false}
			Clean([]string{}, &packages, cache, bf, info, flags)

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).To(BeNil())

			Expect(bytes).To(Equal([]byte(expectedContents)))
		})

		It("Should not modify the existing Brewfile if the --dry-run flag is set", func() {
			flags = Flags{true}
			Clean([]string{}, &packages, cache, bf, info, flags)
			bytes, error := ioutil.ReadFile(bf)
			Expect(error).To(BeNil())

			Expect(bytes).To(Equal([]byte(contents)))
		})
	})

})
