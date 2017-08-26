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

	"sort"

	"regexp"

	"io/ioutil"

	"github.com/boltdb/bolt"
	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	"github.com/spf13/cobra"
)

var addFlags Flags

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&addFlags.DryRun, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")

	addCmd.Flags().BoolVarP(&addFlags.Tap, "tap", "t", false, "add a tap")
	addCmd.Flags().BoolVarP(&addFlags.Brew, "brew", "b", false, "add a brew package")
	addCmd.Flags().BoolVarP(&addFlags.Cask, "cask", "c", false, "add a cask")
	addCmd.Flags().BoolVarP(&addFlags.Mas, "mas", "m", false, "add a mas app")

	addCmd.Flags().StringSliceVar(&addFlags.Args, "args", []string{}, "supply args to be used during installations and updates")
	addCmd.Flags().StringVar(&addFlags.RestartService, "restart-service", "", "always (every time bundle runs), changed (after changes and updates)")
	addCmd.Flags().StringVarP(&addFlags.MasID, "mas-id", "i", "", "id required for mas packages")
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a dependency to your Brewfile",
	Long:  DocsAdd,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var packages brewfile.Packages

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}

		cache := brew.Cache{DB: db}

		err = Add(args, &packages, cache, brewfilePath, addFlags, level)
		errorExit(err)
	},
}

func Add(args []string, packages *brewfile.Packages, cache brew.Cache, brewfilePath string, flags Flags, level int) error {
	if !flagProvided(flags) {
		return ErrNoPackageType("add")
	}

	toAdd := args[0]
	packageType := getPackageType(flags)

	if err := packages.FromBrewfile(brewfilePath); err != nil {
		return err
	}

	if entryExists(string(packages.Bytes()), packageType, toAdd) {
		return ErrEntryAlreadyExists(toAdd)
	}

	cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}

	if err := cacheMap.FromPackages(packages.Brew); err != nil {
		return err
	}

	if err := cacheMap.ResolveDependencyMap(level); err != nil {
		return err
	}

	if flags.Tap {
		if !hasCorrectTapFormat(toAdd) {
			return ErrInvalidTapFormat
		}
		packages.Tap = addPackage(packageType, toAdd, packages.Tap, flags)
		sort.Strings(packages.Tap)
	}

	if flags.Brew {
		updated, err := addBrewPackage(toAdd, flags.RestartService, flags.Args, cacheMap, flags, level)
		if err != nil {
			return err
		}
		packages.Brew = updated
	}

	if flags.Cask {
		packages.Cask = addPackage(packageType, toAdd, packages.Cask, flags)
		sort.Strings(packages.Cask)
	}

	if flags.Mas {
		if !hasMasID(flags.MasID) {
			return ErrNoMasID(toAdd)
		}

		packages.Mas = addPackage(packageType, toAdd, packages.Mas, flags)
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

func addBrewPackage(add, restart string, args []string, cacheMap brew.CacheMap, flags Flags, level int) ([]string, error) {
	if len(restart) > 1 {
		switch restart {
		case "always":
			restart = "true"
		case "changed":
			restart = ":changed"
		default:
			return []string{}, ErrInvalidRestartServiceOption
		}
	}

	if err := cacheMap.Add(brew.Entry{Name: add, RestartService: restart, Args: args}, level); err != nil {
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

func addPackage(packageType, newPackage string, packages []string, flags Flags) []string {
	packageEntry := constructBaseEntry(packageType, newPackage)

	if packageType == "mas" {
		packageEntry = appendMasID(packageEntry, flags.MasID)
	}

	fmt.Printf("Added %s %s to Brewfile.\n", packageType, newPackage)
	return append(packages, packageEntry)
}

func hasCorrectTapFormat(tap string) bool {
	result, _ := regexp.MatchString(`.+/.+`, tap)
	return result
}

func hasMasID(i string) bool {
	return len(i) > 0
}

func appendMasID(packageEntry, i string) string {
	return fmt.Sprintf("%s, id: %s", packageEntry, i)
}
