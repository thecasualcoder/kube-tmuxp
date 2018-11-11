package commander

import (
	"os"
	"os/exec"
)

// Commander is an interface to execute commands
type Commander interface {
	Execute(cmdStr string, args []string, envs []string) (string, error)
}

// Default is a Commander implementation that
// can execute commands on a actual machine
type Default struct{}

// Execute executes a command on the actual machine
func (Default) Execute(cmdStr string, args []string, envs []string) (string, error) {
	cmd := exec.Command(cmdStr, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
