package cmd

import (
	"github.com/lgug2z/bfm/brew"
	"github.com/spf13/cobra"
	"os/exec"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the cache of Homebrew formula information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var cache brew.InfoCache
		infoCommand := exec.Command("brew", "info", "--all", "--json=v1")

		Refresh(args, brewInfoPath, cache, infoCommand)
	},
}

func init() {
	RootCmd.AddCommand(refreshCmd)
}

func Refresh(args []string, brewInfo string, cache brew.InfoCache, command *exec.Cmd) {
	error := cache.Refresh(brewInfo, command)
	errorExit(error)
}
