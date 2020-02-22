package gcloud

import (
	"encoding/json"
	"fmt"
	"github.com/thecasualcoder/kube-tmuxp/pkg/commander"
	"strings"
)

// ListProjects lists the projects for logged-in user
func ListProjects(commander commander.Commander) ([]string, error) {
	args := []string{
		"projects",
		"list",
		"--format=json",
	}
	response, err := commander.Execute("gcloud", args, nil)
	fullCommand := strings.Join(append([]string{"gcloud"}, args...), " ")
	if err != nil {
		return nil, fmt.Errorf("error executing %s: %v", fullCommand, err)
	}
	var parsedResponse []map[string]string
	err = json.Unmarshal([]byte(response), &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling the response from command %s: %v", fullCommand, err)
	}

	projectIds := make([]string, 0, len(parsedResponse))
	for _, project := range parsedResponse {
		projectIds = append(projectIds, project["projectId"])
	}

	return projectIds, nil
}

// Cluster represent the GKE Cluster
type Cluster struct {
	Name      string
	Location  string
	Locations []string
}

func (cluster Cluster) IsRegional() bool {
	return Contains(cluster.Locations, cluster.Location)
}

// Clusters represents the list of Cluster
type Clusters []Cluster

// ListClusters for the given projectId
func ListClusters(cmdr commander.Commander, projectId string) (Clusters, error) {
	args := []string{
		"container",
		"clusters",
		"list",
		"--project",
		projectId,
		"--format=json",
	}
	response, err := cmdr.Execute("gcloud", args, nil)
	fullCommand := strings.Join(append([]string{"gcloud"}, args...), " ")
	if err != nil {
		return nil, fmt.Errorf("error executing %s: %v", fullCommand, err)
	}
	var clusters []Cluster
	err = json.Unmarshal([]byte(response), &clusters)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling the response from command %s: %v", fullCommand, err)
	}

	return clusters, nil
}

// Contains checks if the given string array contains given string
func Contains(items []string, input string) bool {
	for _, item := range items {
		if item == input {
			return true
		}
	}
	return false
}
