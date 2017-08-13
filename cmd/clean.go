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
	"strings"
)

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&d, "dry-run", "d", false, "conduct a dry run without modifying the Brewfile")
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
		contents, err := readFileContents(location)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		lines := strings.Split(contents, "\n")

		tap := getPackages("tap", lines)
		brew := getPackages("brew", lines)
		cask := getPackages("cask", lines)
		mas := getPackages("mas", lines)

		newContents := constructFileContents(tap, brew, cask, mas)

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
