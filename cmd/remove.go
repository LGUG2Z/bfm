// Copyright © 2017 Jade Iqbal <jadeiqbal@fastmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"os"
	"sort"
	"strings"

	"io/ioutil"

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
	removeCmd.Flags().BoolVar(&removeFlags.RemoveAll, "all-unused", false, "remove brew package and its unused required, recommended, optional and build dependencies")
	removeCmd.Flags().BoolVar(&removeFlags.RemovePackageAndRequired, "required-unused", false, "remove brew package and its unused required dependencies")
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a dependency from your Brewfile",
	Long: DocsRemove,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var packages brewfile.Packages

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}

		cache := brew.Cache{DB: db}

		error := Remove(args, &packages, cache, brewfilePath, removeFlags)
		errorExit(error)
	},
}

func Remove(args []string, packages *brewfile.Packages, cache brew.Cache, brewfilePath string, flags Flags) error {
	if !flagProvided(flags) {
		return ErrNoPackageType("remove")
	}

	toRemove := args[0]
	packageType := getPackageType(flags)

	if err := packages.FromBrewfile(brewfilePath); err != nil {
		return err
	}

	if !entryExists(string(packages.Bytes()), packageType, toRemove) {
		return ErrEntryDoesNotExist(toRemove)
	}

	cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}
	cacheMap.FromPackages(packages.Brew)
	cacheMap.ResolveRequiredDependencyMap()

	if flags.Tap {
		packages.Tap = removePackage(packageType, toRemove, packages.Tap)
		sort.Strings(packages.Tap)
	}

	if flags.Brew {
		updated, err := removeBrewPackage(toRemove, cacheMap, flags)
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
		fmt.Println(string(packages.Bytes()))
	} else {
		if err := ioutil.WriteFile(brewfilePath, packages.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}

func removeBrewPackage(remove string, cacheMap brew.CacheMap, flags Flags) ([]string, error) {
	if flags.RemoveAll {
		if err := cacheMap.Remove(remove, brew.RemoveAll); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if flags.RemovePackageAndRequired {
		if err := cacheMap.Remove(remove, brew.RemovePackageAndRequired); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if err := cacheMap.Remove(remove, brew.RemovePackageOnly); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
