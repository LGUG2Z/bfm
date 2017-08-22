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

	"regexp"

	"io/ioutil"

	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	"github.com/spf13/cobra"
	"errors"
)

var addFlags struct {
	brew, tap, cask, mas, dryRun  bool
	args                          []string
	restartService, masID         string
	addPackageAndRequired, addAll bool // user picks which dependencies to also add to the brewfile
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&addFlags.dryRun, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")

	addCmd.Flags().BoolVarP(&addFlags.tap, "tap", "t", false, "add a tap")
	addCmd.Flags().BoolVarP(&addFlags.brew, "brew", "b", false, "add a brew package")
	addCmd.Flags().BoolVarP(&addFlags.cask, "cask", "c", false, "add a cask")
	addCmd.Flags().BoolVarP(&addFlags.mas, "mas", "m", false, "add a mas app")

	addCmd.Flags().StringSliceVar(&addFlags.args, "args", []string{}, "supply args to be used during installations and updates")
	addCmd.Flags().StringVar(&addFlags.restartService, "restart-service", "", "always (every time bundle runs), changed (after changes and updates)")
	addCmd.Flags().StringVarP(&addFlags.masID, "mas-id", "i", "", "id required for mas packages")

	addCmd.Flags().BoolVarP(&addFlags.addPackageAndRequired, "required", "r", false, "add package and all required dependencies")
	addCmd.Flags().BoolVarP(&addFlags.addAll, "all", "a", false, "add package and all required, recommended, optional and build dependencies")
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a dependency to your Brewfile",
	Long: `
Adds the dependency given as an argument to the Brewfile.

This command will modify your Brewfile without creating
a backup. Consider running the command with the --dry-run
flag if using bfm for the first time.

The type must be specified using the appropriate flag.

Taps must conform to the format <user/repo>.

Brew packages can have arguments specified using the --arg
flag (multiple arguments can be separated by using a comma),
and can specify service restart behaviour (always: restart
every time bundle is run, changed: only when updated or
changed) with the --restart-service flag.

MAS apps must specify an id using the --mas-id flag which
can be found by running 'mas search <app>'.

Examples:

bfm add -t homebrew/dupes
bfm add -b vim -a HEAD,with-override-system-vi
bfm add -b crisidev/chunkwm/chunkwm -r changed
bfm add -c macvim
bfm add -m Xcode -i 497799835

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !flagProvided(addFlags.tap, addFlags.brew, addFlags.cask, addFlags.mas) {
			fmt.Println("A package type must be specified. See 'bfm add --help'.")
			os.Exit(1)
		}

		toAdd := args[0]
		packageType := getPackageType(addFlags.tap, addFlags.brew, addFlags.cask, addFlags.mas)

		var packages brewfile.Packages
		error := packages.FromBrewfile(brewfilePath)
		errorExit(error)

		if entryExists(string(packages.Bytes()), packageType, toAdd) {
			fmt.Printf("%s '%s' is already in the Brewfile.\n", packageType, toAdd)
			os.Exit(1)
		}

		var cache brew.InfoCache
		error = cache.Read(brewInfoPath)
		errorExit(error)

		cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}
		cacheMap.FromPackages(packages.Brew)
		cacheMap.ResolveRequiredDependencyMap()

		if addFlags.tap {
			if !hasCorrectTapFormat(packageType) {
				fmt.Printf("Unrecognised tap format. Use the format 'user/repo'.\n")
				os.Exit(1)
			}
			packages.Tap = addPackage(packageType, toAdd, packages.Tap)
			sort.Strings(packages.Tap)
		}

		if addFlags.brew {
			updated, error := addBrewPackage(toAdd, addFlags.restartService, addFlags.args, cacheMap)
			errorExit(error)
			packages.Brew = updated
		}

		if addFlags.cask {
			packages.Cask = addPackage(packageType, toAdd, packages.Cask)
			sort.Strings(packages.Cask)
		}

		if addFlags.mas {
			if !hasMasID(addFlags.masID) {
				fmt.Printf("An id is required for mas apps. Get the id with 'mas search %s and try again.\n", strings.ToLower(toAdd))
				os.Exit(1)
			}

			packages.Mas = addPackage(packageType, toAdd, packages.Mas)
			sort.Strings(packages.Mas)
		}

		if addFlags.dryRun {
			fmt.Println(string(packages.Bytes()))
		} else {
			error := ioutil.WriteFile(brewfilePath, packages.Bytes(), 0644)
			errorExit(error)
		}
	},
}

func addBrewPackage(add, restart string, args []string, cacheMap brew.CacheMap) ([]string, error) {
	if len(restart) > 1 {
		switch restart {
		case "always":
			restart = "true"
		case "changed":
			restart = ":changed"
		default:
			return []string{}, errors.New("Valid options for the --restart-service flag are 'true' and 'changed'.")
		}
	}

	if addFlags.addAll {
		if err := cacheMap.Add(brew.Entry{Name: add, RestartService: restart, Args: args}, brew.AddAll); err != nil {
			return []string{}, err
		}
	} else if addFlags.addPackageAndRequired {
		if err := cacheMap.Add(brew.Entry{Name: add, RestartService: restart, Args: args}, brew.AddPackageAndRequired); err != nil {
			return []string{}, err
		}
	} else {
		if err := cacheMap.Add(brew.Entry{Name: add, RestartService: restart, Args: args}, brew.AddPackageOnly); err != nil {
			return []string{}, err
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

func addPackage(packageType, newPackage string, packages []string) []string {
	packageEntry := constructBaseEntry(packageType, newPackage)

	if packageType == "mas" {
		packageEntry = appendMasID(packageEntry, addFlags.masID)
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
