package brew

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"strings"

	"github.com/boltdb/bolt"
)

type Cache struct {
	DB *bolt.DB
}

// TODO: Figure out a consistent naming convention for casks
func (c *Cache) RefreshCasks(command *exec.Cmd) error {
	b, err := command.Output()
	if err != nil {
		return err
	}

	err = c.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("cask"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	allCasks := strings.Fields(string(b))

	for _, cask := range allCasks {
		err := c.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("cask"))
			if err := b.Put([]byte(cask), []byte(cask)); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) Refresh(command *exec.Cmd) error {
	b, err := command.Output()
	if err != nil {
		return err
	}

	err = c.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("brew"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var allInfo []Info

	if err := json.Unmarshal(b, &allInfo); err != nil {
		return err
	}

	for _, pkg := range allInfo {
		err := c.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("brew"))

			key := pkg.FullName
			value, err := json.Marshal(pkg)
			if err != nil {
				return err
			}

			if err := b.Put([]byte(key), []byte(value)); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (c Cache) FindCask(pkg string) (string, error) {
	var cask string

	err := c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cask"))
		v := b.Get([]byte(pkg))

		if v == nil {
			return ErrCouldNotFindPackageInfo(pkg)
		}

		cask = string(v)

		return nil
	})

	if err != nil {
		return "", err
	}

	return cask, nil
}

func (c Cache) Find(pkg string) (Info, error) {
	var info Info

	err := c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("brew"))
		v := b.Get([]byte(pkg))

		if v == nil {
			return ErrCouldNotFindPackageInfo(pkg)
		}

		err := json.Unmarshal(v, &info)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return Info{}, err
	}

	return info, nil
}
