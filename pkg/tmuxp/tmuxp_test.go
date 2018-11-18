package tmuxp_test

import (
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
	"github.com/arunvelsriram/kube-tmuxp/pkg/tmuxp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tmuxpCfg := tmuxp.New("session", tmuxp.Windows{}, tmuxp.Environment{}, &filesystem.Default{})

	assert.NotNil(t, tmuxpCfg)
}
