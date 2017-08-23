package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	"fmt"
	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Remove", func() {

	var bf string = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")
	var info string = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/test.json")
	//var testInfo string = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testInfo.json")
	var cache brew.InfoCache
	var packages brewfile.Packages

	BeforeEach(func() {
		cache = brew.InfoCache{}
		packages = brewfile.Packages{}
	})

	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {
			error := Remove([]string{"something"}, &brewfile.Packages{}, brew.InfoCache{}, "", "", Flags{})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrNoPackageType("remove").Error()))
		})
	})

	Describe("When the command is called with a flag specifying the package type and a package name", func() {
		It("Should return an error with an explanation if the package is not in the Brewfile", func() {
			Expect(createTestFile(bf, "")).To(Succeed())

			error := Remove([]string{"a2ps"}, &packages, cache, bf, info, Flags{Brew: true})
			Expect(error).To(HaveOccurred())

			errorMessage := error.Error()
			Expect(errorMessage).To(Equal(ErrEntryDoesNotExist("a2ps").Error()))

			Expect(removeTestFile(bf)).To(Succeed())
		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {
			Expect(createTestFile(bf, "brew 'a2ps'")).To(Succeed())

			_ = captureStdout(func() {
				error := Remove([]string{"a2ps"}, &packages, cache, bf, info, Flags{Brew: true, DryRun: true})
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
				error := Remove([]string{"some/repo"}, &packages, cache, bf, info, Flags{Tap: true})
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
				error := Remove([]string{"firefox"}, &packages, cache, bf, info, Flags{Cask: true})
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
				error := Remove([]string{"Xcode"}, &packages, cache, bf, info, Flags{Mas: true})
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

		})
	})

	Describe("When the command is called for a brew entry --required", func() {
		It("Should remove a the brew entry and its required dependencies from the Brewfile", func() {

		})

		It("Should not remove required dependencies that are still required by other packages from the Brewfile", func() {

		})
	})

	Describe("When the command is called for a brew entry --all", func() {
		It("Should remove a the brew entry and its required, recommended and build dependencies from the Brewfile", func() {

		})

		It("Should not remove any dependencies that are still required by other packages from the Brewfile", func() {

		})
	})
})
