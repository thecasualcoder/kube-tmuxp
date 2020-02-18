package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kube-tmuxp",
	Short: `Tool to generate tmuxp configs that help to switch between multiple Kubernetes contexts safely`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configFileName := ".kube-tmuxp.yaml"
	cfgFile = path.Join(home, configFileName)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")
}
