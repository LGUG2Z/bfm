package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"os"
	"bytes"
	"io"
	"fmt"
	"io/ioutil"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func createTestInfoFile(contents string) error {
	file := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testInfo.json")

	error := ioutil.WriteFile(file, []byte(contents), 0644)
	if error != nil {
		return error
	}

	return nil
}

func removeTestInfoFile() error {
	file := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testInfo.json")

	error := os.Remove(file)
	if error != nil {
		return error
	}

	return nil
}

func createTestBrewfile(contents string) error {
	file := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")

	error := ioutil.WriteFile(file, []byte(contents), 0644)
	if error != nil {
		return error
	}

	return nil
}

func removeTestBrewfile() error {
	file := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "/src/github.com/lgug2z/bfm/testData/testBrewfile")

	error := os.Remove(file)
	if error != nil {
		return error
	}

	return nil
}

