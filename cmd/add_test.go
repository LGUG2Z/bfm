package cmd_test

import (
	. "github.com/LGUG2Z/bfm/cmd"

	"fmt"
	"io/ioutil"
	"os"

	"github.com/LGUG2Z/bfm/brew"
	"github.com/LGUG2Z/bfm/brewfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add", func() {
	var (
		bf       = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/LGUG2Z/bfm/testData/testBrewfile")
		dbFile   = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/LGUG2Z/bfm/testData/testDB.bolt")
		cache    brew.Cache
		packages brewfile.Packages
		db       *TestDB
		f        TestFile
	)

	BeforeEach(func() {
		f = TestFile{Path: bf, Contents: ""}
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

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			err := Add([]string{"something"}, &brewfile.Packages{}, cache, "", Flags{}, 0)
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
			t := TestFile{Path: bf + ".customContents", Contents: contents}
			Expect(t.Create()).To(Succeed())
			defer t.Remove()

			err := Add([]string{"a2ps"}, &packages, cache, bf+".customContents", Flags{Brew: true}, 0)
			Expect(err).To(HaveOccurred())

			errorMessage := err.Error()
			Expect(errorMessage).To(Equal(ErrEntryAlreadyExists("a2ps").Error()))
		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			db.AddTestBrewsByName("a2ps")

			_ = captureStdout(func() {
				Expect(Add([]string{"a2ps"}, &packages, cache, bf, Flags{Brew: true, DryRun: true}, brew.Required)).To(Succeed())
			})

			bytes, err := ioutil.ReadFile(bf)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

		})
	})

	Describe("When the command is called for a tap", func() {
		It("Should return an error if the tap format is not user/repo", func() {
			error := Add([]string{"bad:format"}, &brewfile.Packages{}, cache, bf, Flags{Tap: true}, 0)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrInvalidTapFormat.Error()))
		})

		It("Should add a validly formatted tap to the Brewfile", func() {
			_ = captureStdout(func() {
				Expect(Add([]string{"good/format"}, &brewfile.Packages{}, cache, bf, Flags{Tap: true}, 0)).To(Succeed())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("tap 'good/format'\n")))

		})
	})

	Describe("When the command is called for a mas app", func() {
		It("Should return an error if no mas id is provided", func() {
			error := Add([]string{"Xcode"}, &brewfile.Packages{}, cache, bf, Flags{Mas: true}, 0)
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrNoMasID("Xcode").Error()))

		})

		It("Should add a mas app with a mas id to the Brewfile", func() {
			_ = captureStdout(func() {
				Expect(Add([]string{"Xcode"}, &brewfile.Packages{}, cache, bf, Flags{Mas: true, MasID: "123456"}, 0)).To(Succeed())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("mas 'Xcode', id: 123456\n")))

		})
	})

	Describe("When the command is called for a cask", func() {
		It("Should add a cask app to the Brewfile", func() {
			_ = captureStdout(func() {
				Expect(Add([]string{"firefox"}, &brewfile.Packages{}, cache, bf, Flags{Cask: true}, 0)).To(Succeed())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("cask 'firefox'\n")))

		})
	})

	Describe("When called for a brew with the --restart-service flag", func() {
		It("Should return an error explaining the valid options if an invalid option is given", func() {
			_ = captureStdout(func() {
				error := Add([]string{"a2ps"}, &brewfile.Packages{}, cache, bf, Flags{Brew: true, RestartService: "wrong"}, 0)
				Expect(error).To(HaveOccurred())
				Expect(error.Error()).To(Equal(ErrInvalidRestartServiceOption.Error()))
			})

		})

		It("Should add brew with restartService transformed from always to true", func() {
			db.AddTestBrewsByName("a2ps")

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, RestartService: "always"}, 0)).To(Succeed())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: true"))

		})

		It("Should add brew with restartService transform changed to :changed", func() {
			db.AddTestBrewsByName("a2ps")

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, RestartService: "changed"}, 0)).To(Succeed())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: :changed"))

		})
	})

	Describe("When called for a brew with the --args flag", func() {
		It("Should add brew with args ", func() {
			db.AddTestBrewsByName("a2ps")

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true, Args: []string{"one", "two"}}, 0)).To(Succeed())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("args: ['one', 'two']"))

		})
	})

	Describe("When dependency level is set to required", func() {
		It("Should add a brew with its required dependencies to the Brewfile", func() {
			db.AddTestBrewsByName("bash")
			db.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}})

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true}, brew.Required)).To(Succeed())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # [required by: a2ps]"))
		})
	})

	Describe("When dependency level is set to recommended", func() {
		It("Should add a brew with its required and recommended dependencies to the Brewfile", func() {
			db.AddTestBrewsByName("bash", "zsh")
			db.AddTestBrewsFromInfo(brew.Info{FullName: "a2ps", Dependencies: []string{"bash"}, RecommendedDependencies: []string{"zsh"}})

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true}, brew.Recommended)).To(Succeed())

			Expect(packages.Brew).To(HaveLen(3))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # [required by: a2ps]"))
			Expect(packages.Brew[2]).To(Equal("brew 'zsh' # [recommended for: a2ps]"))
		})
	})

	Describe("When dependency level is set to optional", func() {
		It("Should add a brew with its required, recommended and optional dependencies to the Brewfile", func() {
			db.AddTestBrewsByName("bash", "zsh", "fish")
			db.AddTestBrewsFromInfo(brew.Info{
				FullName:                "a2ps",
				Dependencies:            []string{"bash"},
				RecommendedDependencies: []string{"zsh"},
				OptionalDependencies:    []string{"fish"},
			})

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true}, brew.Optional)).To(Succeed())

			Expect(packages.Brew).To(HaveLen(4))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # [required by: a2ps]"))
			Expect(packages.Brew[2]).To(Equal("brew 'fish' # [optional for: a2ps]"))
			Expect(packages.Brew[3]).To(Equal("brew 'zsh' # [recommended for: a2ps]"))
		})

	})

	Describe("When dependency level is set to build", func() {
		It("Should add a brew with its required, recommended, optional and build dependencies to the Brewfile", func() {
			db.AddTestBrewsByName("bash", "zsh", "fish", "sh")
			db.AddTestBrewsFromInfo(brew.Info{
				FullName:                "a2ps",
				Dependencies:            []string{"bash"},
				RecommendedDependencies: []string{"zsh"},
				OptionalDependencies:    []string{"fish"},
				BuildDependencies:       []string{"sh"},
			})

			packages := &brewfile.Packages{}

			Expect(Add([]string{"a2ps"}, packages, cache, bf, Flags{Brew: true}, brew.Build)).To(Succeed())

			Expect(packages.Brew).To(HaveLen(5))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # [required by: a2ps]"))
			Expect(packages.Brew[2]).To(Equal("brew 'fish' # [optional for: a2ps]"))
			Expect(packages.Brew[3]).To(Equal("brew 'sh' # [build for: a2ps]"))
			Expect(packages.Brew[4]).To(Equal("brew 'zsh' # [recommended for: a2ps]"))
		})
	})
})
