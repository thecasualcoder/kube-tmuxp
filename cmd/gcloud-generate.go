package cmd

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
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

var gcloudGenerateCmd = &cobra.Command{
	Use:   "gcloud-generate",
	Short: "Generates configs for kube-tmuxp based on gcloud account",
	Run: func(cmd *cobra.Command, args []string) {
		cmdr := &commander.Default{}
		projects, err := getProjects(cmdr)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			os.Exit(1)
		}
		if !apply {
			printConfigFiles(projects, cmd.OutOrStdout())
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Run with --apply to directly generate tmuxp configs for various Kubernetes contexts\n")
			return
		}
		err = generateKubeTmuxpFiles(cmdr, projects)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), err.Error())
			os.Exit(1)
		}
	},
}

func generateKubeTmuxpFiles(cmdr commander.Commander, projects kubetmuxp.Projects) error {
	fs := &filesystem.Default{}
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

func getProjects(cmdr commander.Commander) (kubetmuxp.Projects, error) {
	gCloudProjects := gcloud.Projects{}
	if projectIDs != nil && len(projectIDs) > 0 {
		for _, projectID := range projectIDs {
			gCloudProjects = append(gCloudProjects, gcloud.Project{ProjectId: projectID})
		}
	} else {
		gCloudProjects = getGCloudProjects(cmdr, allProjects)
	}
	additionalEnvsMap := map[string]string{}
	for _, env := range additionalEnvs {
		envKeyValue := strings.Split(env, "=")
		if len(envKeyValue) != 2 {
			return nil, fmt.Errorf("wrong env format: should be key=value")
		}
		additionalEnvsMap[envKeyValue[0]] = envKeyValue[1]
	}
	projects := make(kubetmuxp.Projects, 0, len(gCloudProjects))
	for _, gCloudProject := range gCloudProjects {
		clusters, err := gcloud.ListClusters(cmdr, gCloudProject.ProjectId)
		if err != nil {
			return nil, err
		}
		_, _ = fmt.Fprintf(os.Stderr, "Number of clusters for %s project: %d\n", gCloudProject, len(clusters))

		kubetmuxpClusters := make(kubetmuxp.Clusters, 0, len(clusters))
		for _, cluster := range clusters {
			zone := ""
			region := ""
			isRegional := cluster.IsRegional()
			if isRegional {
				region = cluster.Location
			} else {
				zone = cluster.Location
			}
			baseEnvs := map[string]string{
				"KUBETMUXP_CLUSTER_NAME":        cluster.Name,
				"KUBETMUXP_CLUSTER_LOCATION":    cluster.Location,
				"KUBETMUXP_CLUSTER_IS_REGIONAL": fmt.Sprintf("%v", isRegional),
				"GCP_PROJECT_ID":                gCloudProject.ProjectId,
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
			Name:     gCloudProject.ProjectId,
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

func getGCloudProjects(cmdr commander.Commander, allProjects bool) gcloud.Projects {
	projects, err := gcloud.ListProjects(cmdr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Number of gcloud projects: %d\n", len(projects))
	if allProjects {
		return projects
	}
	selectedProjects, err := getSelectedProjects(projects)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Number of selected gcloud projects: %d\n", len(selectedProjects))
	return selectedProjects
}

func getSelectedProjects(projects gcloud.Projects) (gcloud.Projects, error) {
	var selectedProjectIDs []string
	prompt := &survey.MultiSelect{
		Message: "Select gcloud projects that you want to configure:",
		Options: projects.IDs(),
		FilterFn: func(s string, options []string) []string {
			var acc []string
			for _, option := range options {
				if fuzzy.Match(s, option) {
					acc = append(acc, option)
				}
			}
			return acc
		},
	}
	opt := func(options *survey.AskOptions) error {
		options.Stdio.Out = os.Stderr
		return nil
	}
	validator := func(ans interface{}) error { return nil }
	err := survey.AskOne(prompt, &selectedProjectIDs, validator, opt)
	if err != nil {
		return nil, fmt.Errorf("error selecting project: %v", err)
	}
	return projects.Filter(selectedProjectIDs), nil
}

var allProjects, apply bool
var additionalEnvs, projectIDs []string

func init() {
	gcloudGenerateCmd.Flags().BoolVar(&allProjects, "allProjects", false, "Skip confirmation for projects")
	gcloudGenerateCmd.Flags().StringSliceVar(&projectIDs, "projectIDs", nil, "Comma separated Project IDs to which the configurations need to be fetched")
	gcloudGenerateCmd.Flags().BoolVar(&apply, "apply", false, "Directly create the tmuxp configs for selected projects")
	gcloudGenerateCmd.Flags().StringSliceVar(&additionalEnvs, "additionalEnvs", nil, "Additional envs to be populated")
	rootCmd.AddCommand(gcloudGenerateCmd)
}
