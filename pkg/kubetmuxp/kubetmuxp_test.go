package kubetmuxp_test

import (
	"io"
	"strings"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/kubetmuxp"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	var reader io.Reader
	cfg := kubetmuxp.NewConfig(reader)

	assert.NotNil(t, cfg)
}

func TestLoadConfig(t *testing.T) {
	t.Run("should load the config", func(t *testing.T) {
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

		cfg := kubetmuxp.NewConfig(reader)
		err := cfg.Load()

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
		assert.Equal(t, cfg.Projects, expectedProjects)
	})

	t.Run("should return error if loading fails", func(t *testing.T) {
		reader := strings.NewReader("invalid yaml")

		cfg := kubetmuxp.NewConfig(reader)
		err := cfg.Load()

		assert.NotNil(t, err)
		assert.Equal(t, cfg.Projects, kubetmuxp.Projects(nil))
	})
}
