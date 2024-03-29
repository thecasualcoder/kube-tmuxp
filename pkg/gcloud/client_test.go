package gcloud

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/kube-tmuxp/pkg/internal/mock"
)

func TestListProjects(t *testing.T) {
	t.Run("should return error if there is an error executing gcloud command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commander := mock.NewCommander(ctrl)
		commander.EXPECT().Execute("gcloud", []string{"projects", "list", "--format=json"}, nil).Return("", fmt.Errorf("please login"))

		projects, err := ListProjects(commander)

		assert.EqualError(t, err, "error executing gcloud projects list --format=json: please login")
		assert.Empty(t, projects)
	})

	t.Run("should return error if there is an unmarshal from gcloud command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commander := mock.NewCommander(ctrl)
		commander.EXPECT().Execute("gcloud", []string{"projects", "list", "--format=json"}, nil).Return("invalid json response", nil)

		projects, err := ListProjects(commander)

		assert.EqualError(t, err, "error unmarshaling the response from command gcloud projects list --format=json: invalid character 'i' looking for beginning of value")
		assert.Empty(t, projects)
	})

	t.Run("should return projects from gcloud command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commander := mock.NewCommander(ctrl)
		commander.EXPECT().Execute("gcloud", []string{"projects", "list", "--format=json"}, nil).Return(`[
  {
    "createTime": "2016-08-20T04:30:54.605Z",
    "lifecycleState": "ACTIVE",
    "name": "My Project",
    "projectId": "clean-pottery",
    "projectNumber": "10",
	"users": [
		"someuser"
	]
  }
]`, nil)

		projects, err := ListProjects(commander)

		assert.NoError(t, err)
		assert.Equal(t, Projects{Project{ProjectId: "clean-pottery"}}, projects)
	})
}

func TestListClusters(t *testing.T) {
	projectId := "projectId"

	t.Run("should return error if there is an error executing gcloud command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commander := mock.NewCommander(ctrl)
		commander.EXPECT().Execute("gcloud", []string{"container",
			"clusters",
			"list",
			"--project",
			projectId,
			"--format=json",
		}, nil).Return("", fmt.Errorf("please login"))

		projects, err := ListClusters(commander, projectId)

		assert.EqualError(t, err, "error executing gcloud container clusters list --project projectId --format=json: please login")
		assert.Empty(t, projects)
	})

	t.Run("should return error if there is an unmarshal from gcloud command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commander := mock.NewCommander(ctrl)
		commander.EXPECT().Execute("gcloud", []string{"container",
			"clusters",
			"list",
			"--project",
			projectId,
			"--format=json",
		}, nil).Return("invalid json response", nil)

		projects, err := ListClusters(commander, projectId)

		assert.EqualError(t, err, "error unmarshaling the response from command gcloud container clusters list --project projectId --format=json: invalid character 'i' looking for beginning of value")
		assert.Empty(t, projects)
	})

	t.Run("should return clusters from gcloud command", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		commander := mock.NewCommander(ctrl)
		commander.EXPECT().Execute("gcloud", []string{"container",
			"clusters",
			"list",
			"--project",
			projectId,
			"--format=json",
		}, nil).Return(`[
  {
    "currentNodeCount": 8,
    "databaseEncryption": {
      "state": "DECRYPTED"
    },
    "defaultMaxPodsConstraint": {
      "maxPodsPerNode": "110"
    },
    "location": "asia-southeast1",
    "locations": [
      "asia-southeast1-a",
      "asia-southeast1-c",
      "asia-southeast1-b"
    ],
    "name": "cluster-one",
    "network": "default",
    "nodeIpv4CidrSize": 24,
    "status": "RUNNING",
    "subnetwork": "default-subnet",
    "zone": "asia-southeast1"
  },
  {
    "currentNodeCount": 8,
    "databaseEncryption": {
      "state": "DECRYPTED"
    },
    "defaultMaxPodsConstraint": {
      "maxPodsPerNode": "110"
    },
    "location": "asia-southeast1",
    "locations": [
      "asia-southeast1-a",
      "asia-southeast1-c",
      "asia-southeast1-b"
    ],
    "name": "cluster-two",
    "network": "default",
    "nodeIpv4CidrSize": 24,
    "status": "RUNNING",
    "subnetwork": "default-subnet",
    "zone": "asia-southeast1"
  }
]`, nil)

		expectedClusters := Clusters{
			Cluster{Name: "cluster-one", Location: "asia-southeast1", Locations: []string{"asia-southeast1-a", "asia-southeast1-c", "asia-southeast1-b"}},
			Cluster{Name: "cluster-two", Location: "asia-southeast1", Locations: []string{"asia-southeast1-a", "asia-southeast1-c", "asia-southeast1-b"}}}

		projects, err := ListClusters(commander, projectId)

		assert.NoError(t, err)
		assert.Equal(t, expectedClusters, projects)
	})
}

func TestProjects_IDs(t *testing.T) {
	projects := Projects{Project{ProjectId: "project_one"}, Project{ProjectId: "project_two"}}

	assert.Equal(t, []string{"project_one", "project_two"}, projects.IDs())
}

func TestProjects_Filter(t *testing.T) {
	projects := Projects{Project{ProjectId: "project_one"}, Project{ProjectId: "project_two"}}

	assert.Equal(t, Projects{Project{ProjectId: "project_one"}}, projects.Filter([]string{"project_one", "invalid_project"}))
}

func TestCluster_IsRegional(t *testing.T) {
	t.Run("should return true for regional clusters", func(t *testing.T) {
		cluster := Cluster{Name: "cluster-one", Location: "asia-southeast1", Locations: []string{"asia-southeast1-a", "asia-southeast1-c", "asia-southeast1-b"}}

		assert.True(t, cluster.IsRegional())
	})

	t.Run("should return false for zonal clusters", func(t *testing.T) {
		cluster := Cluster{Name: "cluster-one", Location: "asia-southeast1-a", Locations: []string{"asia-southeast1-a", "asia-southeast1-c", "asia-southeast1-b"}}

		assert.False(t, cluster.IsRegional())
	})
}
