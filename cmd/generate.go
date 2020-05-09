package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/thecasualcoder/kube-tmuxp/pkg/generator"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/thecasualcoder/kube-tmuxp/pkg/commander"
	"github.com/thecasualcoder/kube-tmuxp/pkg/filesystem"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generates tmuxp configs for various Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		options := generator.Options{
			From:           from,
			AllProjects:    allProjects,
			ProjectIDs:     projectIDs,
			AdditionalEnvs: additionalEnvs,
			Apply:          apply,
			CfgFile:        cfgFile,
		}
		fs := &filesystem.Default{}
		cmdr := &commander.Default{}

		generator, err := generator.NewGenerator(options, fs, cmdr)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			os.Exit(1)
		}
		generator.Generate(cmd.OutOrStderr(), cmd.ErrOrStderr())
	},
}

var cfgFile string
var from string
var allProjects, apply bool
var additionalEnvs, projectIDs []string

func init() {
	generateCmd.Flags().StringVar(&cfgFile, "config", getDefaultConfigPath(), "config file")
	generateCmd.Flags().StringVar(&from, "from", "file", "source from which the tmuxp config files are  generated")
	generateCmd.Flags().BoolVar(&allProjects, "all-projects", false, "Skip confirmation for projects")
	generateCmd.Flags().StringSliceVar(&projectIDs, "project-ids", nil, "Comma separated Project IDs to which the configurations need to be fetched")
	generateCmd.Flags().BoolVar(&apply, "apply", false, "Directly create the tmuxp configs for selected projects")
	generateCmd.Flags().StringSliceVar(&additionalEnvs, "additional-envs", nil, "Additional envs to be populated")
	rootCmd.AddCommand(generateCmd)
}

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configFileName := ".kube-tmuxp.yaml"
	return path.Join(home, configFileName)
}
