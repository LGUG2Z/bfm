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

	"github.com/lgug2z/bfm/brew"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&d, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")

	removeCmd.Flags().BoolVarP(&t, "tap", "t", false, "remove a tap")
	removeCmd.Flags().BoolVarP(&b, "brew", "b", false, "remove a brew package")
	removeCmd.Flags().BoolVarP(&c, "cask", "c", false, "remove a cask")
	removeCmd.Flags().BoolVarP(&m, "mas", "m", false, "remove a mas app")
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
		if !flagProvided(t, b, c, m) {
			fmt.Println("A package type must be specified. See 'bfm remove --help'.")
			os.Exit(1)
		}

		packageToRemove := args[0]
		packageType := getPackageType(t, b, c, m)

		contents, err := readFileContents(location)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !entryExists(contents, packageType, packageToRemove) {
			fmt.Printf("%s '%s' not found in the Brewfile.\n", packageType, packageToRemove)
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

		cacheMap := brew.CacheMap{
			C: &cache,
			M: make(brew.Map),
		}

		cacheMap.FromBrewfile(brewLines)
		cacheMap.ResolveDependencies()

		if packageType == "tap" {
			tapLines = removePackage(packageType, packageToRemove, tapLines)
			sort.Strings(tapLines)
		}

		if packageType == "brew" {
			brewLines, err = removeBrewPackage(packageToRemove, cacheMap)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if packageType == "cask" {
			caskLines = removePackage(packageType, packageToRemove, caskLines)
			sort.Strings(caskLines)
		}

		if packageType == "mas" {
			masLines = removePackage(packageType, packageToRemove, masLines)
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

func removeBrewPackage(remove string, cacheMap brew.CacheMap) ([]string, error) {
	if err := cacheMap.Remove(remove); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := []string{}

	for _, b := range cacheMap.M {
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
