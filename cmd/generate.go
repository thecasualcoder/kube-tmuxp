package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates tmuxp configs for various Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generate")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
