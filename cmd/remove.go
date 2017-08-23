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
	Long: `
Removes the dependency given as an argument from the Brewfile.

This command will modify your Brewfile without creating
a backup. Consider running the command with the --dry-run
flag if using bfm for the first time.

The type must be specified using the appropriate flag.

Examples:

bfm remove -t homebrew/dupes
bfm remove -b vim
bfm remove -c macvim
bfm remove -m Xcode

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !flagProvided(removeFlags) {
			fmt.Println("A package type must be specified. See 'bfm remove --help'.")
			os.Exit(1)
		}

		toRemove := args[0]
		packageType := getPackageType(removeFlags)

		var packages brewfile.Packages
		error := packages.FromBrewfile(brewfilePath)
		errorExit(error)

		if !entryExists(string(packages.Bytes()), packageType, toRemove) {
			fmt.Printf("%s '%s' not found in the Brewfile.\n", packageType, toRemove)
			os.Exit(1)
		}

		var cache brew.InfoCache
		error = cache.Read(brewInfoPath)
		errorExit(error)

		cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}
		cacheMap.FromPackages(packages.Brew)
		cacheMap.ResolveRequiredDependencyMap()

		if removeFlags.Tap {
			packages.Tap = removePackage(packageType, toRemove, packages.Tap)
			sort.Strings(packages.Tap)
		}

		if removeFlags.Brew {
			updated, error := removeBrewPackage(toRemove, cacheMap)
			errorExit(error)
			packages.Brew = updated
		}

		if removeFlags.Cask {
			packages.Cask = removePackage(packageType, toRemove, packages.Cask)
			sort.Strings(packages.Cask)
		}

		if removeFlags.Mas {
			packages.Mas = removePackage(packageType, toRemove, packages.Mas)
			sort.Strings(packages.Mas)
		}

		if removeFlags.DryRun {
			fmt.Println(string(packages.Bytes()))
		} else {
			error := ioutil.WriteFile(brewfilePath, packages.Bytes(), 0644)
			errorExit(error)
		}
	},
}

func removeBrewPackage(remove string, cacheMap brew.CacheMap) ([]string, error) {
	if removeFlags.RemoveAll {
		if err := cacheMap.Remove(remove, brew.RemoveAll); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if removeFlags.RemovePackageAndRequired {
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
