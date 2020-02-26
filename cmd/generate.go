package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/jfreeland/kube-tmuxp/pkg/commander"
	"github.com/jfreeland/kube-tmuxp/pkg/filesystem"
	"github.com/jfreeland/kube-tmuxp/pkg/kubeconfig"
	"github.com/jfreeland/kube-tmuxp/pkg/kubetmuxp"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generates tmuxp configs for various Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		fs := &filesystem.Default{}
		cmdr := &commander.Default{}
		kubeCfg, err := kubeconfig.New(fs, cmdr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Using config file:", cfgFile)
		kubetmuxpCfg, err := kubetmuxp.NewConfig(cfgFile, fs, kubeCfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err = kubetmuxpCfg.Process(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
