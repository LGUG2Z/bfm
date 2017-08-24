package brew_test

import (
	. "github.com/lgug2z/bfm/brew"

	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
)

var _ = Describe("Cache", func() {
	var (
		cache  Cache
		dbFile = fmt.Sprintf("%s/src/github.com/lgug2z/bfm/testData/testDB.bolt", os.Getenv("GOPATH"))
		db     *TestDB
	)

	BeforeEach(func() {
		testDB, err := NewTestDB(dbFile)
		db = testDB
		Expect(err).ToNot(HaveOccurred())
		cache.DB = db.DB
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("With no existing database file", func() {
		It("Should create a new database file and populate with all brew info", func() {
			command := exec.Command("echo", `[ { "full_name": "vim" } ]`)

			Expect(cache.Refresh(command)).To(Succeed())

			actual, err := cache.Find("vim", db.DB)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Info{FullName: "vim"}))
		})
	})

	Describe("With existing database file", func() {
		It("Should update database file and with all new brew info", func() {
			Expect(db.AddTestBrews("vim")).To(Succeed())

			command := exec.Command("echo", `[ { "full_name": "emacs" } ]`)

			Expect(cache.Refresh(command)).To(Succeed())

			actual, err := cache.Find("emacs", db.DB)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Info{FullName: "emacs"}))
		})
	})

	Describe("With a populated Cache", func() {
		It("Should find and return the Info of a package", func() {
			Expect(db.AddTestBrews("vim")).To(Succeed())

			expected := Info{FullName: "vim"}
			actual, err := cache.Find("vim", db.DB)

			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expected))
		})

		It("Should return an error if a package cannot be found", func() {
			Expect(db.AddTestBrews("vim")).To(Succeed())

			_, err := cache.Find("notvim", db.DB)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(ErrCouldNotFindPackageInfo("notvim").Error()))
		})
	})
})
