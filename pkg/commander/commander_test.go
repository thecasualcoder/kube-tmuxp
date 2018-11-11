package commander_test

import (
	"testing"

	"github.com/arunvelsriram/kube-tmuxp/pkg/commander"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Run("should execute command on the machine", func(t *testing.T) {
		cmdr := commander.Default{}
		out, err := cmdr.Execute("echo", []string{"-n", "test"}, []string{})

		assert.Nil(t, err)
		assert.Equal(t, "test", out)
	})

	t.Run("should be able to set env variables for a command", func(t *testing.T) {
		cmdr := commander.Default{}
		out, err := cmdr.Execute("env", []string{}, []string{"TEST_ENV=test"})

		assert.Nil(t, err)
		assert.Contains(t, out, "TEST_ENV=test")
	})

	t.Run("should return error if execution fails", func(t *testing.T) {
		cmdr := commander.Default{}
		out, err := cmdr.Execute("invalid-cmd", []string{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "", out)
	})
}
