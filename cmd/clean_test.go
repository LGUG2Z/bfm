package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	"fmt"
	"os"

	"io/ioutil"

	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clean", func() {
	var (
		bf       = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testBrewfile")
		dbFile   = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")
		cache    brew.Cache
		packages brewfile.Packages
		contents = `
tap 'homebrew/bundle'
brew 'a2ps'
tap 'homebrew/core'
cask 'google-chrome'
mas 'Xcode', id: 497799835
cask 'firefox'
# some comment
`
		f  TestFile
		db *TestDB
	)

	BeforeEach(func() {
		f = TestFile{Path: bf, Contents: contents}
		Expect(f.Create()).To(Succeed())

		testDB, err := NewTestDB(dbFile)
		db = testDB
		Expect(err).ToNot(HaveOccurred())
		cache.DB = db.DB
	})

	AfterEach(func() {
		f.Remove()
		db.Close()
	})

	Describe("When the command is called", func() {
		It("Should read in the packages currently in the Brewfile", func() {
			db.AddTestBrewsByName("a2ps")

			expectedPackages := brewfile.Packages{
				Tap:  []string{"tap 'homebrew/bundle'", "tap 'homebrew/core'"},
				Brew: []string{"brew 'a2ps'"},
				Cask: []string{"cask 'firefox'", "cask 'google-chrome'"},
				Mas:  []string{"mas 'Xcode', id: 497799835"},
			}

			Expect(Clean([]string{}, &packages, cache, bf, Flags{DryRun: false}, 0)).To(Succeed())

			Expect(packages).To(Equal(expectedPackages))
		})

		It("Should not proceed if a package in the Brewfile is not in the BoltDB cache", func() {
			err := Clean([]string{}, &packages, cache, bf, Flags{DryRun: false}, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(brew.ErrCouldNotFindPackageInfo("a2ps").Error()))
		})

		It("Should write out a new Brewfile in alphabetical order split into tap, brew, cask and mas sections", func() {
			expectedContents := `tap 'homebrew/bundle'
tap 'homebrew/core'

brew 'a2ps'

cask 'firefox'
cask 'google-chrome'

mas 'Xcode', id: 497799835
`

			db.AddTestBrewsByName("a2ps")

			Expect(Clean([]string{}, &packages, cache, bf, Flags{DryRun: false}, 0)).To(Succeed())

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).To(BeNil())

			Expect(bytes).To(Equal([]byte(expectedContents)))
		})

		It("Should not modify the existing Brewfile if the --dry-run flag is set", func() {
			db.AddTestBrewsByName("a2ps")

			_ = captureStdout(func() {
				Expect(Clean([]string{}, &packages, cache, bf, Flags{DryRun: true}, 0)).To(Succeed())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).To(BeNil())

			Expect(bytes).To(Equal([]byte(contents)))
		})

		It("Should output the cleaned Brewfile contents to stdout", func() {
			db.AddTestBrewsByName("a2ps")

			expectedOutput := `tap 'homebrew/bundle'
tap 'homebrew/core'

brew 'a2ps'

cask 'firefox'
cask 'google-chrome'

mas 'Xcode', id: 497799835
`

			output := captureStdout(func() {
				Expect(Clean([]string{}, &packages, cache, bf, Flags{DryRun: true}, 0)).To(Succeed())
			})

			Expect(output).To(Equal(expectedOutput))
		})
	})
})
