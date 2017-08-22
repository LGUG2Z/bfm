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

	"io/ioutil"
	"sort"

	"github.com/lgug2z/bfm/brew"
	"github.com/lgug2z/bfm/brewfile"
	"github.com/spf13/cobra"
)

type Flags struct {
	Brew, Tap, Cask, Mas, DryRun        bool
	Args                                []string
	RestartService, MasID               string
	AddPackageAndRequired, AddAll       bool
	RemovePackageAndRequired, RemoveAll bool
}

var cleanFlags Flags

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&cleanFlags.DryRun, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")
}

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean yp your Brewfile",
	Long: `
Cleans up your Brewfile, removing all comments and
sorting all dependencies into alphabetised groups
with the order tap -> brew -> cask -> mas.

This command will modify your Brewfile without creating
a backup. Consider running the command with the --dry-run
flag if using bfm for the first time.

`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var cache brew.InfoCache
		var packages brewfile.Packages
		Clean(args, &packages, cache, brewfilePath, brewInfoPath, cleanFlags)
	},
}

func Clean(args []string, packages *brewfile.Packages, cache brew.InfoCache, brewfilePath, brewInfoPath string, flags Flags) {
	error := packages.FromBrewfile(brewfilePath)
	errorExit(error)

	error = cache.Read(brewInfoPath)
	errorExit(error)

	cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}
	cacheMap.FromPackages(packages.Brew)
	cacheMap.ResolveRequiredDependencyMap()

	cleanBrews, error := cleanBrews(cacheMap)
	errorExit(error)

	packages.Brew = cleanBrews

	if flags.DryRun {
		fmt.Println(string(packages.Bytes()))
	} else {
		if err := ioutil.WriteFile(brewfilePath, packages.Bytes(), 0644); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}
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
