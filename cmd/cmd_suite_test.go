package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
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

func createTestFile(file, contents string) error {
	error := ioutil.WriteFile(file, []byte(contents), 0644)
	if error != nil {
		return error
	}

	return nil
}

func removeTestFile(file string) error {
	error := os.Remove(file)
	if error != nil {
		return error
	}

	return nil
}
