package cmd

import (
	"github.com/boltdb/bolt"
	"github.com/lgug2z/bfm/brew"
	"github.com/spf13/cobra"
	"os/exec"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the cache of Homebrew formula information",
	Long:  DocsRefresh,
	Run: func(cmd *cobra.Command, args []string) {
		brewInfo := exec.Command("brew", "info", "--all", "--json=v1")
		caskInfo := exec.Command("brew", "cask", "search")

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
