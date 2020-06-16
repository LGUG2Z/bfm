package cmd

import (
	"os/exec"

	"github.com/LGUG2Z/bfm/brew"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the cache of brew formula and cask information from tapped repositories",
	Long:  DocsRefresh,
	Run: func(cmd *cobra.Command, args []string) {
		brewInfo := exec.Command("brew", "info", "--all", "--json=v1")
		caskInfo := exec.Command("brew", "search", "--casks")

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				errorExit(err)
			}
		}()

		cache := brew.Cache{DB: db}

		err = Refresh(args, cache, brewInfo, caskInfo)
		errorExit(err)
	},
}

func init() {
	RootCmd.AddCommand(refreshCmd)
}

func Refresh(args []string, cache brew.Cache, brewCommand, caskCommand *exec.Cmd) error {
	if err := cache.Refresh(brewCommand); err != nil {
		return err
	}

	if err := cache.RefreshCasks(caskCommand); err != nil {
		return err
	}

	return nil
}
