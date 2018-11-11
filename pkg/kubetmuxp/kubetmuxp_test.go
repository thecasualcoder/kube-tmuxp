package kubetmuxp_test

import (
	"io"
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

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var reader io.Reader
	kubeCfg := getKubeCfg(ctrl)
	kubetmuxpCfg := kubetmuxp.New(reader, kubeCfg)

	assert.NotNil(t, kubetmuxpCfg)
}

func TestLoadConfig(t *testing.T) {
	t.Run("should load the config", func(t *testing.T) {
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
		kubetmuxpCfg := kubetmuxp.New(reader, kubeCfg)

		err := kubetmuxpCfg.Load()

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

		assert.Nil(t, err)
		assert.Equal(t, kubetmuxpCfg.Projects, expectedProjects)
	})

	t.Run("should return error if loading fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := strings.NewReader("invalid yaml")
		kubeCfg := getKubeCfg(ctrl)
		kubetmuxpCfg := kubetmuxp.New(reader, kubeCfg)

		err := kubetmuxpCfg.Load()

		assert.NotNil(t, err)
		assert.Equal(t, kubetmuxpCfg.Projects, kubetmuxp.Projects(nil))
	})
}

func TestRegional(t *testing.T) {
	t.Run("should return true when region alone is given", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test",
			Region: "test-region",
		}

		regional, err := cluster.Regional()

		assert.Nil(t, err)
		assert.True(t, regional)
	})

	t.Run("should return false when zone alone is given", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name: "test",
			Zone: "test-zone",
		}

		regional, err := cluster.Regional()

		assert.Nil(t, err)
		assert.False(t, regional)
	})

	t.Run("should return error when both region and zone are given", func(t *testing.T) {
		cluster := kubetmuxp.Cluster{
			Name:   "test",
			Region: "test-region",
			Zone:   "test-zone",
		}

		_, err := cluster.Regional()

		assert.EqualError(t, err, "Only one of region or zone should be given")
	})
}
