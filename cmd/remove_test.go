package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Remove", func() {
	Describe("When the command is called without any flags", func() {
		It("Should return an error with info about required flags for specifying package types", func() {

		})
	})

	Describe("When the command is called with a flag specifying the package type and a package name", func() {
		It("Should return an error with an explanation if the package is not in the Brewfile", func() {

		})

		It("Should not modify the Brewfile if the --dry-run flag is set", func() {

		})

		It("Should remove a tap entry from the Brewfile", func() {

		})

		It("Should remove a cask entry from the Brewfile", func() {

		})

		It("Should remove a mas entry from the Brewfile", func() {

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
