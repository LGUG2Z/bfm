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

	"github.com/lgug2z/bfm/brew"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&d, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")

	addCmd.Flags().BoolVarP(&t, "tap", "t", false, "add a tap")
	addCmd.Flags().BoolVarP(&b, "brew", "b", false, "add a brew package")
	addCmd.Flags().BoolVarP(&c, "cask", "c", false, "add a cask")
	addCmd.Flags().BoolVarP(&m, "mas", "m", false, "add a mas app")

	addCmd.Flags().StringSliceVarP(&a, "args", "a", []string{}, "supply args to be used during installations and updates")
	addCmd.Flags().StringVarP(&r, "restart-service", "r", "", "always: every time bundle runs, changed: after changes and updates")
	addCmd.Flags().StringVarP(&i, "mas-id", "i", "", "id required for mas packages")
}

var a []string
var i, r string

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
		if !flagProvided(t, b, c, m) {
			fmt.Println("A package type must be specified. See 'bfm add --help'.")
			os.Exit(1)
		}

		packageToAdd := args[0]
		packageType := getPackageType(t, b, c, m)

		contents, err := readFileContents(location)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if entryExists(contents, packageType, packageToAdd) {
			fmt.Printf("%s '%s' is already in the Brewfile.\n", packageType, packageToAdd)
			os.Exit(1)
		}

		lines := strings.Split(contents, "\n")

		tapLines := getPackages("tap", lines)
		brewLines := getPackages("brew", lines)
		caskLines := getPackages("cask", lines)
		masLines := getPackages("mas", lines)

		var cache brew.InfoCache

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := cache.Read(fmt.Sprintf("%s/%s", home, ".brewInfo.json")); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		brewMap := make(brew.Map)
		brewMap.FromBrewfile(brewLines, &cache)
		brewMap.ResolveDependencies(&cache)

		if packageType == "tap" {
			if !hasCorrectTapFormat(packageType) {
				fmt.Printf("Unrecognised tap format. Use the format 'user/repo'.\n")
				os.Exit(1)
			}
			tapLines = addPackage(packageType, packageToAdd, tapLines)
			sort.Strings(tapLines)
		}

		if packageType == "brew" {
			brewLines, err = addBrewPackage(packageToAdd, r, a, brewMap, &cache)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if packageType == "cask" {
			caskLines = addPackage(packageType, packageToAdd, caskLines)
			sort.Strings(caskLines)
		}

		if packageType == "mas" {
			if !hasMasId(i) {
				fmt.Printf("An id is required for mas apps. Get the id with 'mas search %s and try again.\n", strings.ToLower(packageToAdd))
				os.Exit(1)
			}

			masLines = addPackage(packageType, packageToAdd, masLines)
			sort.Strings(masLines)
		}

		newContents := constructFileContents(tapLines, brewLines, caskLines, masLines)

		if d {
			fmt.Println(newContents)
		} else {
			f, err := os.Create(location)
			if err != nil {
				fmt.Print(err)
			}

			f.WriteString(newContents)
		}
	},
}

func addBrewPackage(add, restart string, args []string, m brew.Map, i *brew.InfoCache) ([]string, error) {
	if err := m.Add(add, restart, args, i); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := []string{}

	for _, b := range m {
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
		packageEntry = appendMasId(packageEntry, i)
	}

	fmt.Printf("Added %s %s to Brewfile.\n", packageType, newPackage)
	return append(packages, packageEntry)
}

func hasCorrectTapFormat(tap string) bool {
	result, _ := regexp.MatchString(`.+/.+`, tap)
	return result
}

func hasMasId(i string) bool {
	return len(i) > 0
}

func appendArgs(packageEntry string, a []string) string {
	withQuotes := []string{}

	for _, arg := range a {
		withQuotes = append(withQuotes, "'"+arg+"'")
	}

	return packageEntry + ", " + "args: [" + strings.Join(withQuotes, ", ") + "]"
}

func appendRestartService(packageEntry, r string) string {
	var option string

	if r == "changed" {
		option = ":changed"
	} else if r == "always" {
		option = "true"
	} else {
		println("Valid options for --restart-services are 'always' and 'changed'.")
		os.Exit(1)
	}

	return fmt.Sprintf("%s, restart_service: %s", packageEntry, option)
}

func appendMasId(packageEntry, i string) string {
	return fmt.Sprintf("%s, id: %s", packageEntry, i)
}
