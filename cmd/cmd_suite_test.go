package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lgug2z/bfm/brew"
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

type TestDB struct {
	*bolt.DB
}

func NewTestDB(path string) (*TestDB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &TestDB{db}, nil
}

func (db *TestDB) Close() {
	defer os.Remove(db.Path())
	db.DB.Close()
}

func (db *TestDB) AddTestBrews(names ...string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("brew"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	for _, n := range names {
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("brew"))

			value, err := json.Marshal(brew.Info{FullName: n})
			if err != nil {
				return err
			}

			if err := b.Put([]byte(n), []byte(value)); err != nil {
				return err
			}

			return nil
		})
	}

	if err != nil {
		return err
	}

	return nil
}

func (db *TestDB) AddTestBrewsFromInfo(info ...brew.Info) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("brew"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	for _, i := range info {
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("brew"))

			value, err := json.Marshal(i)
			if err != nil {
				return err
			}

			if err := b.Put([]byte(i.FullName), []byte(value)); err != nil {
				return err
			}

			return nil
		})
	}

	if err != nil {
		return err
	}

	return nil
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
