package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	"fmt"
	"os"

	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Remove", func() {

	var bf = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")
	var dbFile = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")

	var cache brew.InfoCache
	var packages brewfile.Packages

	BeforeEach(func() {
		cache = brew.InfoCache{}
		packages = brewfile.Packages{}
	})

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			err := Remove([]string{"something"}, &brewfile.Packages{}, brew.InfoCache{}, "", Flags{}, TestDB{}.DB)
			Expect(err).To(HaveOccurred())

			errorMessage := err.Error()
			Expect(errorMessage).To(Equal(ErrNoPackageType("remove").Error()))
		})
	})

	Describe("When the command is called with a flag specifying the package type and a package name", func() {
		It("Should return an error with an explanation if the package is not in the Brewfile", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true}, TestDB{}.DB)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrEntryDoesNotExist("a2ps").Error()))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("a2ps")).To(Succeed())

			Expect(createTestFile(bf, "brew 'a2ps'")).To(Succeed())

			_ = captureStdout(func() {
				error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, DryRun: true}, testDB.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("brew 'a2ps'")))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should remove a tap entry from the Brewfile", func() {
			Expect(createTestFile(bf, "tap 'some/repo'")).To(Succeed())

			_ = captureStdout(func() {
				error := Remove([]string{"some/repo"}, &packages, cache, bf, Flags{Tap: true}, TestDB{}.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should remove a cask entry from the Brewfile", func() {
			Expect(createTestFile(bf, "cask 'firefox'")).To(Succeed())

			_ = captureStdout(func() {
				error := Remove([]string{"firefox"}, &packages, cache, bf, Flags{Cask: true}, TestDB{}.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should remove a mas entry from the Brewfile", func() {
			Expect(createTestFile(bf, "mas 'Xcode', id: 123456")).To(Succeed())

			_ = captureStdout(func() {
				error := Remove([]string{"Xcode"}, &packages, cache, bf, Flags{Mas: true}, TestDB{}.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When the command is called for a brew entry without --required or --all", func() {
		It("Should remove a the brew entry from the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("bash")).To(Succeed())
			Expect(testDB.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})).To(Succeed())

			brewfileContents := `
brew 'a2ps'
brew 'bash' # required by: a2ps
	`
			Expect(createTestFile(bf, brewfileContents)).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, brew.InfoCache{}, bf, Flags{Brew: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(1))
			Expect(packages.Brew[0]).To(Equal("brew 'bash'"))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When the command is called for a brew entry --required", func() {
		It("Should remove a the brew entry and its required dependencies from the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("bash")).To(Succeed())
			Expect(testDB.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})).To(Succeed())

			brewfileContents := `
brew 'a2ps'
brew 'bash' # required by: a2ps
	`
			Expect(createTestFile(bf, brewfileContents)).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, brew.InfoCache{}, bf, Flags{Brew: true, RemovePackageAndRequired: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(0))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should not remove required dependencies that are still required by other packages from the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("bash")).To(Succeed())
			Expect(testDB.AddTestBrewsFromInfo(
				brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}},
				brew.Info{FullName: "zsh", Dependencies: []string{"bash"}},
			)).To(Succeed())

			brewfileContents := `
brew 'a2ps'
brew 'bash' # required by: a2ps, zsh
brew 'zsh'
		`

			Expect(createTestFile(bf, brewfileContents)).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, brew.InfoCache{}, bf, Flags{Brew: true, RemovePackageAndRequired: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'bash' # required by: zsh"))
			Expect(packages.Brew[1]).To(Equal("brew 'zsh'"))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When the command is called for a brew entry --all", func() {
		It("Should remove a the brew entry and its required, recommended and build dependencies from the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("bash", "zsh", "fish", "sh")).To(Succeed())
			Expect(testDB.AddTestBrewsFromInfo(
				brew.Info{
					FullName:                "a2ps",
					Dependencies:            []string{"bash"},
					OptionalDependencies:    []string{"zsh"},
					RecommendedDependencies: []string{"fish"},
					BuildDependencies:       []string{"sh"},
				},
			)).To(Succeed())

			brewfileContents := `
brew 'a2ps'
brew 'bash' # required by: a2ps
brew 'fish'
brew 'sh'
brew 'zsh'
	`

			Expect(createTestFile(bf, brewfileContents)).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, brew.InfoCache{}, bf, Flags{Brew: true, RemoveAll: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(0))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should not remove any dependencies that are still required by other packages from the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("bash", "zsh", "fish", "sh")).To(Succeed())
			Expect(testDB.AddTestBrewsFromInfo(
				brew.Info{
					FullName:                "a2ps",
					Dependencies:            []string{"bash"},
					OptionalDependencies:    []string{"zsh"},
					RecommendedDependencies: []string{"fish"},
					BuildDependencies:       []string{"sh"},
				},
				brew.Info{
					FullName:     "vim",
					Dependencies: []string{"bash"},
				},
			)).To(Succeed())

			brewfileContents := `
brew 'a2ps'
brew 'bash' # required by: a2ps, vim
brew 'fish'
brew 'sh'
brew 'zsh'
brew 'vim'
		`

			Expect(createTestFile(bf, brewfileContents)).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, brew.InfoCache{}, bf, Flags{Brew: true, RemoveAll: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'bash' # required by: vim"))
			Expect(packages.Brew[1]).To(Equal("brew 'vim'"))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})
})
