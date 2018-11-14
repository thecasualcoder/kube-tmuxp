package cmd

import (
	"bufio"
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
		reader, err := os.Open(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bufioReader := bufio.NewReader(reader)

		kubeCfg, err := kubeconfig.New(&filesystem.Default{}, &commander.Default{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		kubetmuxpCfg := kubetmuxp.New(bufioReader, kubeCfg)
		err = kubetmuxpCfg.Load()
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
