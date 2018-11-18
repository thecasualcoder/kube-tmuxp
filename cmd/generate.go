package cmd

import (
	"fmt"
	"os"

	"github.com/arunvelsriram/kube-tmuxp/pkg/commander"
	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
	"github.com/arunvelsriram/kube-tmuxp/pkg/kubeconfig"
	"github.com/arunvelsriram/kube-tmuxp/pkg/kubetmuxp"
	"github.com/spf13/cobra"
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
