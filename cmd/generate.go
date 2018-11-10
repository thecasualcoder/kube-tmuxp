package cmd

import (
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
		kubetmuxpCfg, err := kubetmuxp.NewConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("%+v\n", kubetmuxpCfg)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
