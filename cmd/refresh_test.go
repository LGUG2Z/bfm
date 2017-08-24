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
			command := exec.Command("echo", `[ { "name": "a2ps", "full_name": "a2ps" } ]`)
			dbFile := fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/lgug2z/bfm/testData/testDB.bolt")

			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()

			Refresh([]string{}, brew.InfoCache{}, command, testDB.DB)

			var v []byte
			err = testDB.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("brew"))
				v = b.Get([]byte("a2ps"))
				return nil
			})

			Expect(v).ToNot(BeNil())

			var info brew.Info
			Expect(json.Unmarshal(v, &info)).To(Succeed())

			Expect(info).To(Equal(brew.Info{Name: "a2ps", FullName: "a2ps"}))
		})
	})
})
