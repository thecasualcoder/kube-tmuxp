package gcloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thecasualcoder/kube-tmuxp/pkg/commander"
)

// ListProjects lists the projects for logged-in user
func ListProjects(commander commander.Commander) (Projects, error) {
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
	var projects []Project
	err = json.Unmarshal([]byte(response), &projects)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling the response from command %s: %v", fullCommand, err)
	}

	return projects, nil
}

// Project represent the GCP project
type Project struct {
	ProjectId string `json:"projectId"`
}

// Projects represent the list of GCP projects
type Projects []Project

func (p Projects) IDs() []string {
	acc := make([]string, 0, len(p))
	for _, project := range p {
		acc = append(acc, project.ProjectId)
	}
	return acc
}

func (p Projects) Filter(projectIDs []string) Projects {
	projectMap := map[string]Project{}
	for _, project := range p {
		projectMap[project.ProjectId] = project
	}
	result := make(Projects, 0, len(p))
	for _, id := range projectIDs {
		if project, ok := projectMap[id]; ok {
			result = append(result, project)
		}
	}
	return result
}

// Cluster represent the GKE Cluster
type Cluster struct {
	Name      string
	Location  string
	Locations []string
}

func (cluster Cluster) IsRegional() bool {
	return !Contains(cluster.Locations, cluster.Location)
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
