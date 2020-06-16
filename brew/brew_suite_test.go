package brew_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"testing"

	"encoding/json"
	"fmt"

	"github.com/LGUG2Z/bfm/brew"
	"github.com/boltdb/bolt"
)

func TestBrew(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brew Suite")
}

var testPath = fmt.Sprintf("%s/%s", os.Getenv("GOPATH"), "src/github.com/LGUG2Z/bfm/testData")

var _ = BeforeSuite(func() {
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		Expect(os.Mkdir(testPath, os.ModePerm)).To(Succeed())
	}
})

var _ = AfterSuite(func() {
	if _, err := os.Stat(testPath); err == nil {
		Expect(os.Remove(testPath)).To(Succeed())
	}
})

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
