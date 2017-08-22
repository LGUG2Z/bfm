package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	"fmt"
	"github.com/lgug2z/bfm/brew"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var _ = Describe("Refresh", func() {
	Describe("When the command is called", func() {
		It("It should populate the file from the output of the given command and write it to disk", func() {
			command := exec.Command("echo", `[ { "name": "a2ps", "full_name": "a2ps", "desc": "Any-to-PostScript filter" } ]`)
			testFile := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/refreshTest.json")
			infoCache := brew.InfoCache{}

			Refresh([]string{}, testFile, infoCache, command)

			bytes, error := ioutil.ReadFile(testFile)
			Expect(error).To(BeNil())
			Expect(strings.TrimSpace(string(bytes))).To(Equal(`[ { "name": "a2ps", "full_name": "a2ps", "desc": "Any-to-PostScript filter" } ]`))

			error = os.Remove(testFile)
			Expect(error).To(BeNil())
		})
	})

})
