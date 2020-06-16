package cmd

import (
	"fmt"

	"sort"

	"github.com/LGUG2Z/bfm/brew"
	"github.com/LGUG2Z/bfm/brewfile"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var cleanFlags Flags

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&cleanFlags.DryRun, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")
}

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up your Brewfile",
	Long:  DocsClean,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var packages brewfile.Packages

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}

		cache := brew.Cache{DB: db}

		err = Clean(args, &packages, cache, brewfilePath, cleanFlags, level)
		errorExit(err)
	},
}

func Clean(args []string, packages *brewfile.Packages, cache brew.Cache, brewfilePath string, flags Flags, level int) error {
	if err := packages.FromBrewfile(brewfilePath); err != nil {
		return err
	}

	cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}
	if err := cacheMap.FromPackages(packages.Brew); err != nil {
		return err
	}

	if err := cacheMap.ResolveDependencyMap(level); err != nil {
		return err
	}

	cleanBrews, err := cleanBrews(cacheMap)
	if err != nil {
		return err
	}

	packages.Brew = cleanBrews

	if flags.DryRun {
		b, err := packages.Bytes()
		if err != nil {
			return err
		}

		fmt.Printf(string(b))
	} else {
		if err := writeToFile(brewfilePath, packages); err != nil {
			return err
		}
	}

	return nil
}

func cleanBrews(cacheMap brew.CacheMap) ([]string, error) {
	clean := []string{}

	for _, b := range cacheMap.Map {
		entry, err := b.Format()
		if err != nil {
			return []string{}, err
		}

		clean = append(clean, entry)
	}

	sort.Strings(clean)
	return clean, nil
}
