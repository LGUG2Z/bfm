package cmd

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyExists               = func(name string) error { return fmt.Errorf("Entry for %s already exists in the Brewfile.", name) }
	ErrInvalidTapFormat            = errors.New("Invalid tap format. See bfm add --help.")
	ErrInvalidRestartServiceOption = errors.New("Invalid --restart-service option. See bfm add --help")
	ErrNoPackageType               = errors.New("No package type specified. See bfm add --help.")
	ErrNoMasID                     = func(name string) error {
		return fmt.Errorf("An ID is required for mas entries. Run 'mas search %s' to get the ID.", name)
	}
)
