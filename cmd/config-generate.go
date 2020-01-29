package cmd

import (
	"fmt"
	"github.com/thecasualcoder/kube-tmuxp/pkg/commander"
	"github.com/thecasualcoder/kube-tmuxp/pkg/filesystem"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubeconfig"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"

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
		projects, err := getProjects()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			os.Exit(1)
		}
		if !apply {
			printConfigFiles(projects, cmd.OutOrStdout())
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Run with --apply to directly generate tmuxp configs for various Kubernetes contexts\n")
			return
		}
		err = generateKubeTmuxpFiles(projects)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			os.Exit(1)
		}
	},
}

func generateKubeTmuxpFiles(projects kubetmuxp.Projects) error {
	fs := &filesystem.Default{}
	cmdr := &commander.Default{}
	kubeCfg, err := kubeconfig.New(fs, cmdr)

	config, err := kubetmuxp.NewConfigWithProjects(projects, fs, kubeCfg)
	if err != nil {
		return err
	}
	return config.Process()
}

func printConfigFiles(projects kubetmuxp.Projects, outStream io.Writer) {
	bytes, err := yaml.Marshal(map[string]kubetmuxp.Projects{"projects": projects})
	if err != nil {
		_, _ = fmt.Fprintln(outStream, err)
		os.Exit(1)
	}
	fmt.Println(string(bytes))
}

func getProjects() (kubetmuxp.Projects, error) {
	projectIds := getGCloudProjects(allProjects)
	additionalEnvsMap := map[string]string{}
	for _, env := range additionalEnvs {
		envKeyValue := strings.Split(env, "=")
		if len(envKeyValue) != 2 {
			return nil, fmt.Errorf("wrong env format: should be key=value")
		}
		additionalEnvsMap[envKeyValue[0]] = envKeyValue[1]
	}
	projects := make(kubetmuxp.Projects, 0, len(projectIds))
	for _, projectId := range projectIds {
		clusters, err := gcloud.ListClusters(projectId)
		if err != nil {
			return nil, err
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
			baseEnvs := map[string]string{
				"KUBETMUXP_CLUSTER_NAME":        cluster.Name,
				"KUBETMUXP_CLUSTER_LOCATION":    cluster.Location,
				"KUBETMUXP_CLUSTER_IS_REGIONAL": fmt.Sprintf("%v", cluster.IsRegional),
				"GCP_PROJECT_ID":                projectId,
			}
			kubetmuxpClusters = append(kubetmuxpClusters, kubetmuxp.Cluster{
				Name:    cluster.Name,
				Zone:    zone,
				Region:  region,
				Context: cluster.Name,
				Envs:    mergeEnvs(baseEnvs, additionalEnvsMap),
			})
		}
		projects = append(projects, kubetmuxp.Project{
			Name:     projectId,
			Clusters: kubetmuxpClusters,
		})
	}
	return projects, nil
}

func mergeEnvs(base, additionalEnvsMap map[string]string) map[string]string {
	for k, v := range additionalEnvsMap {
		expandedValue := os.Expand(v, func(s string) string {
			value, ok := base[s]
			if ok {
				return value
			} else {
				return os.Getenv(s)
			}
		})
		base[k] = expandedValue
	}
	return base
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
	opt := func(options *survey.AskOptions) error {
		options.Stdio.Out = os.Stderr
		return nil
	}
	err := survey.AskOne(prompt, &selectedProjects, func(ans interface{}) error { return nil }, opt)
	return selectedProjects, err
}

var allProjects, apply bool
var additionalEnvs []string

func init() {
	configGenerateCmd.Flags().BoolVar(&allProjects, "allProjects", false, "Skip confirmation for projects")
	configGenerateCmd.Flags().BoolVar(&apply, "apply", false, "Directly create the tmuxp configs for selected projects")
	configGenerateCmd.Flags().StringSliceVar(&additionalEnvs, "additionalEnvs", nil, "Additional envs to be populated")
	rootCmd.AddCommand(configGenerateCmd)
}
