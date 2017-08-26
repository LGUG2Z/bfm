package brew

import "fmt"

const (
	RequiredDependency = iota
	RecommendedDependency
	OptionalDependency
	BuildDependency
)

const (
	Required = iota
	Recommended
	Optional
	Build
)

var (
	ErrCouldNotFindPackageInfo = func(name string) error {
		return fmt.Errorf("Could not find information for %s. Aborting.\n"+
			"If this package is from a new tap, run 'bfm refresh' to use info from the new tap.\n"+
			"With manually added taps the full name format should be used: 'github_user/repo/package'.\n", name)
	}
)
