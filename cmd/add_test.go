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

var _ = Describe("Add", func() {
	var (
		bf       = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testBrewfile")
		dbFile   = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")
		cache    brew.Cache
		packages brewfile.Packages
	)

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			err := Add([]string{"something"}, &brewfile.Packages{}, cache, "", Flags{})
			Expect(err).To(HaveOccurred())

			errorMessage := err.Error()
			Expect(errorMessage).To(Equal(ErrNoPackageType("add").Error()))
		})
	})

	Describe("When the command is called with a flag specifying the package type and a package name", func() {
		It("Should return an error with an explanation if the package is already in the Brewfile", func() {
			contents := `
tap 'homebrew/bundle'
brew 'a2ps'
tap 'homebrew/core'
cask 'google-chrome'
mas 'Xcode', id: 497799835
cask 'firefox'
# some comment
	`
			t := TestFile{Path: bf, Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			err := Add([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true})
			Expect(err).To(HaveOccurred())

			errorMessage := err.Error()
			Expect(errorMessage).To(Equal(ErrEntryAlreadyExists("a2ps").Error()))
		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("a2ps")

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				err := Add([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, DryRun: true})
				Expect(err).ToNot(HaveOccurred())
			})

			bytes, err := ioutil.ReadFile(bf)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

		})
	})

	Describe("When the command is called for a tap", func() {
		It("Should return an error if the tap format is not user/repo", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Add([]string{"bad:format"}, &brewfile.Packages{}, cache, bf, Flags{Tap: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrInvalidTapFormat.Error()))
		})

		It("Should add a validly formatted tap to the Brewfile", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Add([]string{"good/format"}, &brewfile.Packages{}, cache, bf, Flags{Tap: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("tap 'good/format'\n")))

		})
	})

	Describe("When the command is called for a mas app", func() {
		It("Should return an error if no mas id is provided", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			error := Add([]string{"Xcode"}, &brewfile.Packages{}, cache, bf, Flags{Mas: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrNoMasID("Xcode").Error()))

		})

		It("Should add a mas app with a mas id to the Brewfile", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Add([]string{"Xcode"}, &brewfile.Packages{}, cache, bf, Flags{Mas: true, MasID: "123456"})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("mas 'Xcode', id: 123456")))

		})
	})

	Describe("When the command is called for a cask", func() {
		It("Should add a cask app to the Brewfile", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Add([]string{"firefox"}, &brewfile.Packages{}, cache, bf, Flags{Cask: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("cask 'firefox'\n")))

		})
	})

	Describe("When called for a brew with the --restart-service flag", func() {
		It("Should return an error explaining the valid options if an invalid option is given", func() {
			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			_ = captureStdout(func() {
				error := Add([]string{"a2ps"}, &brewfile.Packages{}, cache, bf, Flags{Brew: true, RestartService: "wrong"})
				Expect(error).To(HaveOccurred())
				Expect(error.Error()).To(Equal(ErrInvalidRestartServiceOption.Error()))
			})

		})

		It("Should add brew with restartService transformed from always to true", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("a2ps")

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, RestartService: "always"})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: true"))

		})

		It("Should add brew with restartService transform changed to :changed", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("a2ps")

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, RestartService: "changed"})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: :changed"))

		})
	})

	Describe("When called for a brew with the --args flag", func() {
		It("Should add brew with args ", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("a2ps")

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, Args: []string{"one", "two"}})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("args: ['one', 'two']"))

		})
	})

	Describe("When called for a brew without --required or --all", func() {
		It("Should add brew without any of its dependencies to the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("bash")
			testDB.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			packages := &brewfile.Packages{}
			error := Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(1))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))

		})
	})

	Describe("When called for a with --required", func() {
		It("Should add brew and all of its required dependencies to the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("bash")
			testDB.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, AddPackageAndRequired: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # required by: a2ps"))

		})
	})

	Describe("When called for a brew with --all", func() {
		It("Should add brew with all of its required, recommended, optional and build dependencies to the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			testDB.AddTestBrewsByName("bash", "zsh", "fish", "sh")
			testDB.AddTestBrewsFromInfo(brew.Info{
				FullName:                "a2ps",
				Dependencies:            []string{"bash"},
				OptionalDependencies:    []string{"zsh"},
				RecommendedDependencies: []string{"fish"},
				BuildDependencies:       []string{"sh"},
			})

			t := TestFile{Path: bf, Contents: ""}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, AddAll: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(5))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # required by: a2ps"))
			Expect(packages.Brew[2]).To(Equal("brew 'fish'"))
			Expect(packages.Brew[3]).To(Equal("brew 'sh'"))
			Expect(packages.Brew[4]).To(Equal("brew 'zsh'"))

		})
	})
})
