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

	"github.com/spf13/cobra"
)

var checkFlags Flags

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolVarP(&checkFlags.Tap, "tap", "t", false, "check a tap")
	checkCmd.Flags().BoolVarP(&checkFlags.Brew, "brew", "b", false, "check a brew package")
	checkCmd.Flags().BoolVarP(&checkFlags.Cask, "cask", "c", false, "check a cask")
	checkCmd.Flags().BoolVarP(&checkFlags.Mas, "mas", "m", false, "check a mas app")
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if a dependency is in your Brewfile",
	Long:  DocsCheck,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !flagProvided(checkFlags) {
			fmt.Println("A package type must be specified. See 'bfm check --help'.")
			os.Exit(1)
		}

		packageToCheck := args[0]
		packageType := getPackageType(checkFlags)

		contents, err := readFileContents(brewfilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if entryExists(contents, packageType, packageToCheck) {
			fmt.Printf("%s '%s' is in the Brewfile.\n", packageType, packageToCheck)
		} else {
			fmt.Printf("%s '%s' is not in the Brewfile.\n", packageType, packageToCheck)
		}
	},
}
