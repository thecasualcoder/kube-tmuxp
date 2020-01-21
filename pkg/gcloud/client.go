package gcloud

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"google.golang.org/api/cloudresourcemanager/v1"
	container2 "google.golang.org/genproto/googleapis/container/v1"
)

// ListProjects lists the projects for logged-in user
func ListProjects() ([]string, error) {
	newService, err := cloudresourcemanager.NewService(context.TODO())
	if err != nil {
		return nil, err
	}
	projectsService := cloudresourcemanager.NewProjectsService(newService)

	response, err := projectsService.List().Do()
	if err != nil {
		return nil, err
	}

	projectIds := make([]string, 0, len(response.Projects))
	for _, project := range response.Projects {
		projectIds = append(projectIds, project.ProjectId)
	}

	return projectIds, nil
}

// Cluster represent the GKE Cluster
type Cluster struct {
	Name       string
	Location   string
	IsRegional bool
}

// Clusters represents the list of Cluster
type Clusters []Cluster

// ListClusters for the given projectId
func ListClusters(projectId string) (Clusters, error) {
	clusterManagerClient, err := container.NewClusterManagerClient(context.TODO())
	if err != nil {
		return nil, err
	}

	listClustersRequest := container2.ListClustersRequest{Parent: fmt.Sprintf("projects/%s/locations/%s", projectId, "-")}
	response, err := clusterManagerClient.ListClusters(context.TODO(), &listClustersRequest)
	if err != nil {
		return nil, err
	}
	clusters := make(Clusters, 0, len(response.Clusters))
	for _, cluster := range response.Clusters {
		clusters = append(clusters, Cluster{Name: cluster.Name, Location: cluster.Location, IsRegional: !Contains(cluster.Locations, cluster.Location)})
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
