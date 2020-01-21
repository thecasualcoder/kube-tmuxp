package cmd

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"

	"github.com/spf13/cobra"
	"github.com/thecasualcoder/kube-tmuxp/pkg/gcloud"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubetmuxp"
	"gopkg.in/AlecAivazis/survey.v1"
)

var configGenerateCmd = &cobra.Command{
	Use:     "config-generate",
	Aliases: []string{"cgen"},
	Short:   "Generates configs for kube-tmuxp based on gcloud account",
	Run: func(cmd *cobra.Command, args []string) {

		projectIds := getGCloudProjects(allProjects)

		projects := make(kubetmuxp.Projects, 0, len(projectIds))
		for _, projectId := range projectIds {
			clusters, err := gcloud.ListClusters(projectId)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			_, _ = fmt.Fprintf(os.Stderr, "Number of clusters for %s project: %d\n", projectId, len(clusters))

			kubetmuxpClusters := make(kubetmuxp.Clusters, 0, len(clusters))
			for _, cluster := range clusters {
				zone := ""
				region := ""
				if cluster.IsRegional {
					region = cluster.Location
				} else {
					zone = cluster.Location
				}
				kubetmuxpClusters = append(kubetmuxpClusters, kubetmuxp.Cluster{
					Name:    cluster.Name,
					Zone:    zone,
					Region:  region,
					Context: cluster.Name,
					Envs:    nil,
				})
			}
			projects = append(projects, kubetmuxp.Project{
				Name:     projectId,
				Clusters: kubetmuxpClusters,
			})
		}
		bytes, err := yaml.Marshal(map[string]kubetmuxp.Projects{"projects": projects})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	},
}

func getGCloudProjects(allProjects bool) []string {
	projectIds, err := gcloud.ListProjects()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Number of gcloud projects: %d\n", len(projectIds))
	if allProjects {
		return projectIds
	}
	selectedProjectIds, err := getSelectedProjects(projectIds)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Number of selected gcloud projects: %d\n", len(selectedProjectIds))
	return selectedProjectIds
}

func getSelectedProjects(projects []string) ([]string, error) {
	var selectedProjects []string
	prompt := &survey.MultiSelect{
		Message: "Select gcloud projects that you want to configure:",
		Options: projects,
	}
	err := survey.AskOne(prompt, &selectedProjects, func(ans interface{}) error { return nil })
	return selectedProjects, err
}

var allProjects bool

func init() {
	configGenerateCmd.PersistentFlags().BoolVar(&allProjects, "allProjects", false, "Skip confirmation for projects")
	rootCmd.AddCommand(configGenerateCmd)
}
