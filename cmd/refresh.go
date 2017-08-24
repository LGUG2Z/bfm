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
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var cache brew.InfoCache
		infoCommand := exec.Command("brew", "info", "--all", "--json=v1")

		db, err := bolt.Open(boltFilePath, 0600, nil)
		if err != nil {
			errorExit(err)
		}
		defer db.Close()

		err = Refresh(args, cache, infoCommand, db)
		errorExit(err)
	},
}

func init() {
	RootCmd.AddCommand(refreshCmd)
}

func Refresh(args []string, cache brew.InfoCache, command *exec.Cmd, db *bolt.DB) error {
	if err := cache.Refresh(db, command); err != nil {
		return err
	}

	return nil
}
