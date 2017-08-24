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
		infoCommand := exec.Command("brew", "info", "--all", "--json=v1")

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}
		defer db.Close()

		cache := brew.Cache{DB: db}

		err = Refresh(args, cache, infoCommand)
		errorExit(err)
	},
}

func init() {
	RootCmd.AddCommand(refreshCmd)
}

func Refresh(args []string, cache brew.Cache, command *exec.Cmd) error {
	if err := cache.Refresh(command); err != nil {
		return err
	}

	return nil
}
