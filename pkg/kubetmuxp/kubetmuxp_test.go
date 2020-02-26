package kubetmuxp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/kube-tmuxp/pkg/internal/mock"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubeconfig"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubetmuxp"
)

func getKubeCfg(ctrl *gomock.Controller) kubeconfig.KubeConfig {
	mockCmdr := mock.NewCommander(ctrl)
	mockFS := mock.NewFileSystem(ctrl)
	mockFS.EXPECT().HomeDir().Return("/Users/test", nil)
	kubeCfg, _ := kubeconfig.New(mockFS, mockCmdr)
	return kubeCfg
}

func TestNewConfig(t *testing.T) {
	t.Run("should create a new kube-tmuxp config from file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		reader := strings.NewReader("")
		mockFS.EXPECT().Open("kube-tmuxp-config.yaml").Return(reader, nil)

		kubetmuxpCfg, err := kubetmuxp.NewConfig("kube-tmuxp-config.yaml", mockFS, kubeconfig.KubeConfig{})

		assert.Nil(t, err)
		assert.NotNil(t, kubetmuxpCfg)
	})

	t.Run("should return error if file cannot be opened", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		mockFS.EXPECT().Open("kube-tmuxp-config.yaml").Return(nil, fmt.Errorf("some error"))

		_, err := kubetmuxp.NewConfig("kube-tmuxp-config.yaml", mockFS, kubeconfig.KubeConfig{})

		assert.EqualError(t, err, "some error")
	})

	t.Run("should read the kube-tmuxp configs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		content := `
projects:
- name: test-project
  clusters:
  - name: test-cluster
    zone: test-zone
    region: test-region
    context: test-ctx
    tmux_envs:
      TEST_ENV: test-value`
		reader := strings.NewReader(content)
		mockFS.EXPECT().Open("kube-tmuxp-config.yaml").Return(reader, nil)
		kubetmuxpCfg, _ := kubetmuxp.NewConfig("kube-tmuxp-config.yaml", mockFS, kubeconfig.KubeConfig{})

		expectedProjects := kubetmuxp.Projects{
			{
				Name: "test-project",
				Clusters: kubetmuxp.Clusters{
					{
						Name:    "test-cluster",
						Zone:    "test-zone",
						Region:  "test-region",
						Context: "test-ctx",
						TmuxEnvs: kubetmuxp.TmuxEnvs{
							"TEST_ENV": "test-value",
						},
					},
				},
			},
		}
		assert.NotNil(t, kubetmuxpCfg)
		assert.Equal(t, expectedProjects, kubetmuxpCfg.Projects)
	})

	t.Run("should return error if config cannot be read", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFS := mock.NewFileSystem(ctrl)
		reader := strings.NewReader("invalid yaml")
		mockFS.EXPECT().Open("kube-tmuxp-config.yaml").Return(reader, nil)
		_, err := kubetmuxp.NewConfig("kube-tmuxp-config.yaml", mockFS, kubeconfig.KubeConfig{})

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

		name, err := cluster.DefaultGKEContextName("test-project")

		assert.Nil(t, err)
		assert.Equal(t, "gke_test-project_test-region_test-cluster", name)
	})

	t.Run("should return default context name for zonal cluster", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name: "test-cluster",
			Zone: "test-zone",
		}

		name, err := cluster.DefaultGKEContextName("test-project")

		assert.Nil(t, err)
		assert.Equal(t, "gke_test-project_test-zone_test-cluster", name)
	})

	t.Run("should return error if cluster type (regional or zonal) cannot be determined", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test-cluster",
			Region: "test-region",
			Zone:   "test-zone",
		}

		_, err := cluster.DefaultGKEContextName("test-project")

		assert.EqualError(t, err, "Only one of region or zone should be given")
	})
}
