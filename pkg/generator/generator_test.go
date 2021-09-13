package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/kube-tmuxp/pkg/file"
	"github.com/thecasualcoder/kube-tmuxp/pkg/gcloud"
)

func TestNewGenerator(t *testing.T) {

	t.Run("should fail if invalid from source is given", func(t *testing.T) {
		generator, err := NewGenerator(Options{From: "invalid"}, nil, nil)

		assert.EqualError(t, err, "invalid source provided: valid sources are file,gcloud")
		assert.Nil(t, generator)
	})

	t.Run("should create gcloud generator for gcloud option", func(t *testing.T) {
		generator, err := NewGenerator(Options{From: "gcloud"}, nil, nil)

		assert.Nil(t, err)
		assert.Equal(t, generator, gcloud.NewGenerator(nil, false, nil, false))
	})

	t.Run("should create file generator if options are valid", func(t *testing.T) {
		generator, err := NewGenerator(Options{From: "file"}, nil, nil)

		assert.Nil(t, err)
		assert.Equal(t, generator, file.NewGenerator(nil, nil, ""))
	})

	t.Run("should fail if if options are invalid for file generator", func(t *testing.T) {
		generator, err := NewGenerator(Options{From: "file", AllProjects: true, ProjectIDs: []string{"project1"}, AdditionalEnvs: []string{"a=1"}}, nil, nil)

		assert.Nil(t, generator)
		assert.EqualError(t, err, "error in the flags for source type 'file': \n 1) project-ids should be empty for source file\n 2) all-projects should be false for source file\n 3) additional-envs should be empty for source file\n")
	})
}
