package cmd

import (
	"fmt"

	"sort"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	"github.com/spf13/cobra"
)

var removeFlags Flags

func init() {
	RootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&removeFlags.DryRun, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")

	removeCmd.Flags().BoolVarP(&removeFlags.Tap, "tap", "t", false, "remove a tap")
	removeCmd.Flags().BoolVarP(&removeFlags.Brew, "brew", "b", false, "remove a brew package")
	removeCmd.Flags().BoolVarP(&removeFlags.Cask, "cask", "c", false, "remove a cask")
	removeCmd.Flags().BoolVarP(&removeFlags.Mas, "mas", "m", false, "remove a mas app")
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a dependency from your Brewfile",
	Long:  DocsRemove,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var packages brewfile.Packages

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}

		cache := brew.Cache{DB: db}

		error := Remove(args, &packages, cache, brewfilePath, removeFlags, level)
		errorExit(error)
	},
}

func Remove(args []string, packages *brewfile.Packages, cache brew.Cache, brewfilePath string, flags Flags, level int) error {
	if !flagProvided(flags) {
		return ErrNoPackageType("remove")
	}

	toRemove := args[0]
	packageType := getPackageType(flags)

	if err := packages.FromBrewfile(brewfilePath); err != nil {
		return err
	}

	b, err := packages.Bytes()
	if err != nil {
		return err
	}

	if !entryExists(string(b), packageType, toRemove) {
		return ErrEntryDoesNotExist(toRemove)
	}

	cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}

	if err := cacheMap.FromPackages(packages.Brew); err != nil {
		return err
	}

	if err := cacheMap.ResolveDependencyMap(level); err != nil {
		return err
	}

	if flags.Tap {
		packages.Tap = removePackage(packageType, toRemove, packages.Tap)
		sort.Strings(packages.Tap)
	}

	if flags.Brew {
		updated, err := removeBrewPackage(toRemove, cacheMap, level)
		if err != nil {
			return err
		}
		packages.Brew = updated
	}

	if flags.Cask {
		packages.Cask = removePackage(packageType, toRemove, packages.Cask)
		sort.Strings(packages.Cask)
	}

	if flags.Mas {
		packages.Mas = removePackage(packageType, toRemove, packages.Mas)
		sort.Strings(packages.Mas)
	}

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

func removeBrewPackage(remove string, cacheMap brew.CacheMap, level int) ([]string, error) {
	if err := cacheMap.Remove(remove, level); err != nil {
		return []string{}, err
	}

	lines := []string{}

	for _, b := range cacheMap.Map {
		entry, err := b.Format()
		if err != nil {
			return []string{}, err
		}

		lines = append(lines, entry)
	}

	sort.Strings(lines)
	return lines, nil
}

func removePackage(packageType, packageToRemove string, packages []string) []string {
	updatedPackages := []string{}
	entryToRemove := constructBaseEntry(packageType, packageToRemove)

	for _, p := range packages {
		if !strings.HasPrefix(p, entryToRemove) {
			updatedPackages = append(updatedPackages, p)
		} else {
			fmt.Printf("Removed %s '%s' from Brewfile.\n", packageType, packageToRemove)
		}
	}

	return updatedPackages
}
