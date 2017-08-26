package cmd

import (
	"errors"
	"fmt"
)

var (
	ErrEntryAlreadyExists          = func(name string) error { return fmt.Errorf("Entry for %s already exists in the Brewfile.", name) }
	ErrEntryDoesNotExist           = func(name string) error { return fmt.Errorf("Entry for %s does not exist in the Brewfile.", name) }
	ErrInvalidTapFormat            = errors.New("Invalid tap format. See bfm add --help.")
	ErrInvalidRestartServiceOption = errors.New("Invalid --restart-service option. See bfm add --help")
	ErrDependencyLevelNotSet       = errors.New("BFM_LEVEL not set in shell rc file. See bfm --help.")
	ErrBrewfileNotSet              = errors.New("BFM_BREWFILE not set in shell rc file. See bfm --help.")

	ErrNoPackageType = func(command string) error {
		return fmt.Errorf("No package type specified. See bfm %s --help.", command)
	}
	ErrNoMasID = func(name string) error {
		return fmt.Errorf("An ID is required for mas entries. Run 'mas search %s' to get the ID.", name)
	}
)
