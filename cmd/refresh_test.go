package cmd_test

import (
	. "github.com/lgug2z/bfm/cmd"

	"fmt"
	"os"
	"os/exec"

	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/lgug2z/bfm/brew"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Refresh", func() {
	Describe("When the command is called", func() {
		It("It should populate the file from the output of the given command and write it to disk", func() {
			brewCommand := exec.Command("echo", `[ { "name": "a2ps", "full_name": "a2ps" } ]`)
			caskCommand := exec.Command("echo", `firefox    google-chrome   opera`)
			dbFile := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")

			db, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer db.Close()
			cache := brew.Cache{DB: db.DB}

			Expect(Refresh([]string{}, cache, brewCommand, caskCommand)).To(Succeed())

			var info brew.Info
			var opera, firefox, chrome []byte
			err = db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("brew"))
				v := b.Get([]byte("a2ps"))

				Expect(v).ToNot(BeNil())
				Expect(json.Unmarshal(v, &info)).To(Succeed())

				c := tx.Bucket([]byte("cask"))
				opera = c.Get([]byte("opera"))
				Expect(opera).ToNot(BeNil())
				firefox = c.Get([]byte("firefox"))
				Expect(firefox).ToNot(BeNil())
				chrome = c.Get([]byte("google-chrome"))
				Expect(chrome).ToNot(BeNil())

				return nil
			})

			Expect(err).To(BeNil())

			Expect(info).To(Equal(brew.Info{Name: "a2ps", FullName: "a2ps"}))
			Expect(string(opera)).To(Equal("opera"))
			Expect(string(firefox)).To(Equal("firefox"))
			Expect(string(chrome)).To(Equal("google-chrome"))
		})
	})
})
