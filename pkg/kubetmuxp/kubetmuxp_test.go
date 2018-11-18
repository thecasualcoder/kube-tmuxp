package kubetmuxp_test

import (
	"strings"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/internal/mock"
	"github.com/arunvelsriram/kube-tmuxp/pkg/kubeconfig"
	"github.com/arunvelsriram/kube-tmuxp/pkg/kubetmuxp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getKubeCfg(ctrl *gomock.Controller) kubeconfig.KubeConfig {
	mockCmdr := mock.NewCommander(ctrl)
	mockFS := mock.NewFileSystem(ctrl)
	mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
	kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
	return kubeCfg
}

func TestNewConfig(t *testing.T) {
	t.Run("should create a new kube-tmuxp config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := strings.NewReader("")
		kubeCfg := getKubeCfg(ctrl)
		kubetmuxpCfg, err := kubetmuxp.NewConfig(reader, kubeCfg)

		assert.Nil(t, err)
		assert.NotNil(t, kubetmuxpCfg)
	})

	t.Run("should read the kube-tmuxp configs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		content := `
projects:
- name: test-project
  clusters:
  - name: test-cluster
    zone: test-zone
    region: test-region
    context: test-ctx
    envs:
      TEST_ENV: test-value`
		reader := strings.NewReader(content)
		kubeCfg := getKubeCfg(ctrl)
		kubetmuxpCfg, _ := kubetmuxp.NewConfig(reader, kubeCfg)

		assert.NotNil(t, kubetmuxpCfg)

		expectedProjects := kubetmuxp.Projects{
			{
				Name: "test-project",
				Clusters: kubetmuxp.Clusters{
					{
						Name:    "test-cluster",
						Zone:    "test-zone",
						Region:  "test-region",
						Context: "test-ctx",
						Envs: kubetmuxp.Envs{
							"TEST_ENV": "test-value",
						},
					},
				},
			},
		}
		assert.Equal(t, expectedProjects, kubetmuxpCfg.Projects)
	})

	t.Run("should return error if config cannot be read", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := strings.NewReader("invalid yaml")
		kubeCfg := getKubeCfg(ctrl)
		_, err := kubetmuxp.NewConfig(reader, kubeCfg)

		assert.NotNil(t, err)
	})
}

func TestIsRegional(t *testing.T) {
	t.Run("should return true when region alone is given", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test",
			Region: "test-region",
		}

		regional, err := cluster.IsRegional()

		assert.Nil(t, err)
		assert.True(t, regional)
	})

	t.Run("should return false when zone alone is given", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name: "test",
			Zone: "test-zone",
		}

		regional, err := cluster.IsRegional()

		assert.Nil(t, err)
		assert.False(t, regional)
	})

	t.Run("should return error when both region and zone are given", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test",
			Region: "test-region",
			Zone:   "test-zone",
		}

		_, err := cluster.IsRegional()

		assert.EqualError(t, err, "Only one of region or zone should be given")
	})
}

func TestGeneratedContextName(t *testing.T) {
	t.Run("should return default context name for regional cluster", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test-cluster",
			Region: "test-region",
		}

		name, err := cluster.DefaultContextName("test-project")

		assert.Nil(t, err)
		assert.Equal(t, "gke_test-project_test-region_test-cluster", name)
	})

	t.Run("should return default context name for zonal cluster", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name: "test-cluster",
			Zone: "test-zone",
		}

		name, err := cluster.DefaultContextName("test-project")

		assert.Nil(t, err)
		assert.Equal(t, "gke_test-project_test-zone_test-cluster", name)
	})

	t.Run("should return error if cluster type (regional or zonal) cannot be determined", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test-cluster",
			Region: "test-region",
			Zone:   "test-zone",
		}

		_, err := cluster.DefaultContextName("test-project")

		assert.EqualError(t, err, "Only one of region or zone should be given")
	})
}
