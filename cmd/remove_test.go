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

	var (
		bf       = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")
		dbFile   = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")
		cache    brew.Cache
		packages brewfile.Packages
		db       *TestDB
	)

	BeforeEach(func() {
		testDB, err := NewTestDB(dbFile)
		db = testDB
		Expect(err).ToNot(HaveOccurred())
		cache.DB = db.DB
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			err := Remove([]string{"something"}, &brewfile.Packages{}, cache, "", Flags{})
			Expect(err).To(HaveOccurred())

			errorMessage := err.Error()
			Expect(errorMessage).To(Equal(ErrNoPackageType("remove").Error()))
		})
	})

	Describe("When the command is called with a flag specifying the package type and a package name", func() {
		It("Should return an error with an explanation if the package is not in the Brewfile", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrEntryDoesNotExist("a2ps").Error()))

		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			Expect(db.AddTestBrewsByName("a2ps")).To(Succeed())

			t := TestFile{Path: bf, Contents: "brew 'a2ps'"}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, DryRun: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("brew 'a2ps'")))

		})

		It("Should remove a tap entry from the Brewfile", func() {
			t := TestFile{Path: bf, Contents: "tap 'some/repo'"}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Remove([]string{"some/repo"}, &packages, cache, bf, Flags{Tap: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

		})

		It("Should remove a cask entry from the Brewfile", func() {
			t := TestFile{Path: bf, Contents: "cask 'firefox'"}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Remove([]string{"firefox"}, &packages, cache, bf, Flags{Cask: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

		})

		It("Should remove a mas entry from the Brewfile", func() {
			t := TestFile{Path: bf, Contents: "mas 'Xcode', id: 123456"}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Remove([]string{"Xcode"}, &packages, cache, bf, Flags{Mas: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

		})
	})

	Describe("When the command is called for a brew entry without --required or --all", func() {
		It("Should remove a the brew entry from the Brewfile", func() {
			Expect(db.AddTestBrewsByName("bash")).To(Succeed())
			Expect(db.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})).To(Succeed())

			contents := `
brew 'a2ps'
brew 'bash' # required by: a2ps
	`
			t := TestFile{Path: bf, Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(1))
			Expect(packages.Brew[0]).To(Equal("brew 'bash'"))

		})
	})

	Describe("When the command is called for a brew entry --required", func() {
		It("Should remove a the brew entry and its required dependencies from the Brewfile", func() {
			Expect(db.AddTestBrewsByName("bash")).To(Succeed())
			Expect(db.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})).To(Succeed())

			contents := `
brew 'a2ps'
brew 'bash' # required by: a2ps
	`
			t := TestFile{Path: bf, Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, RemovePackageAndRequired: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(0))

		})

		It("Should not remove required dependencies that are still required by other packages from the Brewfile", func() {
			Expect(db.AddTestBrewsByName("bash")).To(Succeed())
			Expect(db.AddTestBrewsFromInfo(
				brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}},
				brew.Info{FullName: "zsh", Dependencies: []string{"bash"}},
			)).To(Succeed())

			contents := `
brew 'a2ps'
brew 'bash' # required by: a2ps, zsh
brew 'zsh'
		`

			t := TestFile{Path: bf, Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, RemovePackageAndRequired: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'bash' # required by: zsh"))
			Expect(packages.Brew[1]).To(Equal("brew 'zsh'"))

		})
	})

	Describe("When the command is called for a brew entry --all", func() {
		It("Should remove a the brew entry and its required, recommended and build dependencies from the Brewfile", func() {
			Expect(db.AddTestBrewsByName("bash", "zsh", "fish", "sh")).To(Succeed())
			Expect(db.AddTestBrewsFromInfo(
				brew.Info{
					FullName:                "a2ps",
					Dependencies:            []string{"bash"},
					OptionalDependencies:    []string{"zsh"},
					RecommendedDependencies: []string{"fish"},
					BuildDependencies:       []string{"sh"},
				},
			)).To(Succeed())

			contents := `
brew 'a2ps'
brew 'bash' # required by: a2ps
brew 'fish'
brew 'sh'
brew 'zsh'
	`

			t := TestFile{Path: bf, Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, RemoveAll: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(0))

		})

		It("Should not remove any dependencies that are still required by other packages from the Brewfile", func() {
			Expect(db.AddTestBrewsByName("bash", "zsh", "fish", "sh")).To(Succeed())
			Expect(db.AddTestBrewsFromInfo(
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

			contents := `
brew 'a2ps'
brew 'bash' # required by: a2ps, vim
brew 'fish'
brew 'sh'
brew 'zsh'
brew 'vim'
		`

			t := TestFile{Path: bf, Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Remove([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, RemoveAll: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'bash' # required by: vim"))
			Expect(packages.Brew[1]).To(Equal("brew 'vim'"))

		})
	})
})
