package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildVersion string

// SetVersion set the major and minor version
func SetVersion(version string) {
	buildVersion = version
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(buildVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
