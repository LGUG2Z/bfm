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
	var bf string = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")
	var info string = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/test.json")
	var testInfo string = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testInfo.json")

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			error := Add([]string{"something"}, &brewfile.Packages{}, brew.InfoCache{}, "", "", Flags{})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal("A package type must be specified. See 'bfm add -help'."))
		})
	})

	Describe("When the command is called with a flag specifying the package type and a package name", func() {
		var (
			cache    brew.InfoCache
			packages brewfile.Packages
		)

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
			Expect(createTestBrewfile(contents)).To(Succeed())

			error := Add([]string{"a2ps"}, &packages, cache, bf, info, Flags{Brew: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal("brew 'a2ps' is already in the Brewfile."))

			Expect(removeTestBrewfile()).To(Succeed())
		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"a2ps"}, &packages, cache, bf, info, Flags{Brew: true, DryRun: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("")))

			Expect(removeTestBrewfile()).To(Succeed())
		})
	})

	Describe("When the command is called for a tap", func() {
		It("Should return an error if the tap format is not user/repo", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			error := Add([]string{"bad:format"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Tap: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal("Unrecognised tap format. Use the format 'user/repo'."))

			Expect(removeTestBrewfile()).To(Succeed())
		})

		It("Should add a validly formatted tap to the Brewfile", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"good/format"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Tap: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("tap 'good/format'\n")))

			Expect(removeTestBrewfile()).To(Succeed())
		})
	})

	Describe("When the command is called for a mas app", func() {
		It("Should return an error if no mas id is provided", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			error := Add([]string{"Xcode"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Mas: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal("An id is required for mas apps. Get the id with 'mas search Xcode' and try again."))

			Expect(removeTestBrewfile()).To(Succeed())
		})

		It("Should add a mas app with a mas id to the Brewfile", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"Xcode"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Mas: true, MasID: "123456"})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("mas 'Xcode', id: 123456")))

			Expect(removeTestBrewfile()).To(Succeed())
		})
	})

	Describe("When the command is called for a cask", func() {
		It("Should add a cask app to the Brewfile", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"firefox"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Cask: true})
				Expect(error).ToNot(HaveOccurred())
			})

			bytes, error := ioutil.ReadFile(bf)
			Expect(error).ToNot(HaveOccurred())
			Expect(bytes).To(Equal([]byte("cask 'firefox'\n")))

			Expect(removeTestBrewfile()).To(Succeed())
		})
	})

	Describe("When called for a brew with the --restart-service flag", func() {
		It("Should return an error explaining the valid options if an invalid option is given", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			_ = captureStdout(func() {
				error := Add([]string{"a2ps"}, &brewfile.Packages{}, brew.InfoCache{}, bf, info, Flags{Brew: true, RestartService: "wrong"})
				Expect(error).To(HaveOccurred())
				Expect(error.Error()).To(Equal("Valid options for the --restart-service flag are 'true' and 'changed'."))
			})

			Expect(removeTestBrewfile()).To(Succeed())
		})

		It("Should add brew with restartService transformed from always to true", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, info, Flags{Brew: true, RestartService: "always"})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: true"))

			Expect(removeTestBrewfile()).To(Succeed())
		})

		It("Should add brew with restartService transform changed to :changed", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, info, Flags{Brew: true, RestartService: "changed"})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("restart_service: :changed"))

			Expect(removeTestBrewfile()).To(Succeed())
		})
	})

	Describe("When called for a brew with the --args flag", func() {
		It("Should add brew with args ", func() {
			Expect(createTestBrewfile("")).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, info, Flags{Brew: true, Args: []string{"one", "two"}})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).ToNot(BeEmpty())
			Expect(packages.Brew[0]).To(ContainSubstring("args: ['one', 'two']"))

			Expect(removeTestBrewfile()).To(Succeed())
		})
	})

	Describe("When called for a brew without --required or --all", func() {
		It("Should add brew without any of its dependencies to the Brewfile", func() {
			contents := `
[
	{
		"name": "a2ps",
		"full_name": "a2ps",
		"desc": "Any-to-PostScript filter",
		"dependencies": ["bash"]
	},
	{
		"name": "bash" ,
		"full_name": "bash"
	}
]`
			Expect(createTestBrewfile("")).To(Succeed())
			Expect(createTestInfoFile(contents)).To(Succeed())

			packages := &brewfile.Packages{}
			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, testInfo, Flags{Brew: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(1))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))

			Expect(removeTestBrewfile()).To(Succeed())
			Expect(removeTestInfoFile()).To(Succeed())
		})
	})

	Describe("When called for a with --required", func() {
		It("Should add brew and all of its required dependencies to the Brewfile", func() {
			contents := `
[
	{
		"name": "a2ps",
		"full_name": "a2ps",
		"desc": "Any-to-PostScript filter",
		"dependencies": ["bash"]
	},
	{
		"name": "bash" ,
		"full_name": "bash"
	}
]`

			Expect(createTestBrewfile("")).To(Succeed())
			Expect(createTestInfoFile(contents)).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, testInfo, Flags{Brew: true, AddPackageAndRequired: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(2))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # required by: a2ps"))

			Expect(removeTestBrewfile()).To(Succeed())
			Expect(removeTestInfoFile()).To(Succeed())
		})
	})

	Describe("When called for a brew with --all", func() {
		It("Should add brew with all of its required, recommended, optional and build dependencies to the Brewfile", func() {
			contents := `
[
	{
		"name": "a2ps",
		"full_name": "a2ps",
		"desc": "Any-to-PostScript filter",
		"dependencies": ["bash"],
		"optional_dependencies": ["zsh"],
		"recommended_dependencies": ["fish"],
		"build_dependencies": ["sh"]
	},
	{ "name": "bash" , "full_name": "bash" },
	{ "name": "zsh" , "full_name": "zsh" },
	{ "name": "sh" , "full_name": "sh" },
	{ "name": "fish" , "full_name": "fish" }
]`

			Expect(createTestBrewfile("")).To(Succeed())
			Expect(createTestInfoFile(contents)).To(Succeed())

			packages := &brewfile.Packages{}

			error := Add([]string{"a2ps"}, packages, brew.InfoCache{}, bf, testInfo, Flags{Brew: true, AddAll: true})
			Expect(error).ToNot(HaveOccurred())

			Expect(packages.Brew).To(HaveLen(5))
			Expect(packages.Brew[0]).To(Equal("brew 'a2ps'"))
			Expect(packages.Brew[1]).To(Equal("brew 'bash' # required by: a2ps"))
			Expect(packages.Brew[2]).To(Equal("brew 'fish'"))
			Expect(packages.Brew[3]).To(Equal("brew 'sh'"))
			Expect(packages.Brew[4]).To(Equal("brew 'zsh'"))

			Expect(removeTestBrewfile()).To(Succeed())
			Expect(removeTestInfoFile()).To(Succeed())
		})
	})
})
