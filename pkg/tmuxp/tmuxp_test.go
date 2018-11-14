package tmuxp_test

import (
	"io"
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/tmuxp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var writer io.Writer
	tmuxpCfg := tmuxp.New(writer)

	assert.NotNil(t, tmuxpCfg)
}
