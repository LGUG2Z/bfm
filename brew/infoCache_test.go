package brew_test

import (
	. "github.com/lgug2z/bfm/brew"

	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
)

var _ = Describe("InfoCache", func() {
	var (
		cache  InfoCache
		dbFile = fmt.Sprintf("%s/src/github.com/lgug2z/bfm/testData/testDB.bolt", os.Getenv("GOPATH"))
	)

	Describe("With no existing database file", func() {
		It("Should create a new database file and populate with all brew info", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			command := exec.Command("echo", `[ { "full_name": "vim" } ]`)

			err = cache.Refresh(testDB.DB, command)
			Expect(err).ToNot(HaveOccurred())

			actual, err := cache.Find("vim", testDB.DB)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Info{FullName: "vim"}))
		})
	})

	Describe("With existing database file", func() {
		It("Should update database file and with all new brew info", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("vim")).To(Succeed())

			command := exec.Command("echo", `[ { "full_name": "emacs" } ]`)

			err = cache.Refresh(testDB.DB, command)
			Expect(err).ToNot(HaveOccurred())

			actual, err := cache.Find("emacs", testDB.DB)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Info{FullName: "emacs"}))
		})
	})

	Describe("With a populated InfoCache", func() {
		It("Should find and return the Info of a package", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("vim")).To(Succeed())

			expected := Info{FullName: "vim"}
			actual, error := cache.Find("vim", testDB.DB)

			Expect(error).To(BeNil())
			Expect(actual).To(Equal(expected))
		})

		It("Should return an error if a package cannot be found", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Expect(testDB.AddTestBrews("vim")).To(Succeed())

			_, error := cache.Find("notvim", testDB.DB)

			Expect(error).ToNot(BeNil())
		})
	})
})
