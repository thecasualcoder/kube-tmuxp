package cmd

import (
	"bufio"
	"fmt"
	"os"

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

		kubetmuxpCfg := kubetmuxp.NewConfig(bufio.NewReader(reader))
		err = kubetmuxpCfg.Load()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		kubetmuxpCfg.Process()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
