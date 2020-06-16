package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/LGUG2Z/bfm/brew"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bfm",
	Short: "Manage the contents of your Brewfile.",
	Long:  DocsRoot,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(boltPath); os.IsNotExist(err) {
			fmt.Printf("Cache not found. Building...")
			refreshCmd.Run(refreshCmd, []string{""})
			fmt.Printf(" Done.\n\n")
		}
	}}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	brewfilePath string
	boltPath     string
	level        int
)

func init() {
	cobra.OnInitialize(initConfig)
}

type Flags struct {
	Brew, Tap, Cask, Mas, DryRun bool
	Args                         []string
	RestartService, MasID        string
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bfm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bfm")

		viper.SetEnvPrefix("bfm")
		viper.AutomaticEnv() // read in environment variables that match
		brewfilePath = viper.GetString("brewfile")

		if len(brewfilePath) < 1 {
			fmt.Println(ErrBrewfileNotSet)
			os.Exit(1)
		}

		level, err = resolveDependencyLevel(viper.GetString("level"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		boltPath = fmt.Sprintf("%s/%s", home, ".bfm.bolt")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func resolveDependencyLevel(level string) (int, error) {
	switch strings.ToLower(level) {
	case "required":
		return brew.Required, nil
	case "recommended":
		return brew.Recommended, nil
	case "optional":
		return brew.Optional, nil
	case "build":
		return brew.Build, nil
	}

	return 0, ErrDependencyLevelNotSet
}
