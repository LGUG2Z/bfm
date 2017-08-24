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
		info     = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/test.json")
		testInfo = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testInfo.json")
		dbFile   = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")
		cache    brew.InfoCache
		packages brewfile.Packages
	)

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			error := Add([]string{"something"}, &brewfile.Packages{}, brew.InfoCache{}, "", "", Flags{}, TestDB{}.DB)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
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
			Expect(createTestFile(bf, contents)).To(Succeed())

			error := Add([]string{"a2ps"}, &packages, cache, bf, info, Flags{Brew: true}, TestDB{}.DB)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrEntryAlreadyExists("a2ps").Error()))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("a2ps")

			Expect(createTestFile(bf, "")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"a2ps"}, &packages, cache, bf, info, Flags{Brew: true, DryRun: true}, testDB.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When the command is called for a tap", func() {
		It("Should return an error if the tap format is not user/repo", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			error := Add([]string{"bad:format"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Tap: true}, TestDB{}.DB)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrInvalidTapFormat.Error()))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should add a validly formatted tap to the Brewfile", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"good/format"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Tap: true}, TestDB{}.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("tap 'good/format'\n")))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When the command is called for a mas app", func() {
		It("Should return an error if no mas id is provided", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			error := Add([]string{"Xcode"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Mas: true}, TestDB{}.DB)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrNoMasID("Xcode").Error()))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should add a mas app with a mas id to the Brewfile", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"Xcode"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Mas: true, MasID: "123456"}, TestDB{}.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("mas 'Xcode', id: 123456")))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When the command is called for a cask", func() {
		It("Should add a cask app to the Brewfile", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"firefox"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Cask: true}, TestDB{}.DB)
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("cask 'firefox'\n")))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When called for a brew with the --restart-service flag", func() {
		It("Should return an error explaining the valid options if an invalid option is given", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"a2ps"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Brew: true, RestartService: "wrong"}, TestDB{}.DB)
				Expect(error).To(HaveOccurred())
				Expect(error.Error()).To(Equal(ErrInvalidRestartServiceOption.Error()))
			})

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should add brew with restartService transformed from always to true", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("a2ps")

			Expect(createTestFile(bf, "")).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, info, Flags{Brew: true, RestartService: "always"}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: true"))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should add brew with restartService transform changed to :changed", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("a2ps")

			Expect(createTestFile(bf, "")).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, info, Flags{Brew: true, RestartService: "changed"}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: :changed"))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When called for a brew with the --args flag", func() {
		It("Should add brew with args ", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("a2ps")

			Expect(createTestFile(bf, "")).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, info, Flags{Brew: true, Args: []string{"one", "two"}}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("args: ['one', 'two']"))

			Expect(removeTestFile(bf)).To(Succeed())
		})
	})

	Describe("When called for a brew without --required or --all", func() {
		It("Should add brew without any of its dependencies to the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("bash")
			testDB.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})

			contents := `[]`
			Expect(createTestFile(bf, "")).To(Succeed())
			Expect(createTestFile(testInfo, contents)).To(Succeed())

			packages := &brewfile.Packages{}
			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, testInfo, Flags{Brew: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(1))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))

			Expect(removeTestFile(bf)).To(Succeed())
			Expect(removeTestFile(testInfo)).To(Succeed())
		})
	})

	Describe("When called for a with --required", func() {
		It("Should add brew and all of its required dependencies to the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("bash")
			testDB.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})

			contents := `[]`
			Expect(createTestFile(bf, "")).To(Succeed())
			Expect(createTestFile(testInfo, contents)).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, testInfo, Flags{Brew: true, AddPackageAndRequired: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # required by: a2ps"))

			Expect(removeTestFile(bf)).To(Succeed())
			Expect(removeTestFile(testInfo)).To(Succeed())
		})
	})

	Describe("When called for a brew with --all", func() {
		It("Should add brew with all of its required, recommended, optional and build dependencies to the Brewfile", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			testDB.AddTestBrews("bash", "zsh", "fish", "sh")
			testDB.AddTestBrewsFromInfo(brew.Info{
				FullName:                "a2ps",
				Dependencies:            []string{"bash"},
				OptionalDependencies:    []string{"zsh"},
				RecommendedDependencies: []string{"fish"},
				BuildDependencies:       []string{"sh"},
			})

			contents := `[]`
			Expect(createTestFile(bf, "")).To(Succeed())
			Expect(createTestFile(testInfo, contents)).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, testInfo, Flags{Brew: true, AddAll: true}, testDB.DB)
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(5))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # required by: a2ps"))
			Expect(packages.Brew[2]).To(Equal("brew 'fish'"))
			Expect(packages.Brew[3]).To(Equal("brew 'sh'"))
			Expect(packages.Brew[4]).To(Equal("brew 'zsh'"))

			Expect(removeTestFile(bf)).To(Succeed())
			Expect(removeTestFile(testInfo)).To(Succeed())
		})
	})
})
