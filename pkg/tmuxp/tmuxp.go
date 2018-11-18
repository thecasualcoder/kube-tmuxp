package tmuxp

import (
	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
)

// Window represents a window in a tmux session
type Window struct {
	Name  string
	Panes []struct{}
}

// Windows represents a list of windows in a tmux session
type Windows []Window

// Environment represents env variables to be loaded in a tmux session
type Environment map[string]string

// Config represents tmuxp config
type Config struct {
	SessionName string
	Windows
	Environment
	filesystem filesystem.FileSystem
}

// New returns a new tmuxp config
func New(sessionName string, windows Windows, environment Environment, fs filesystem.FileSystem) Config {
	return Config{
		SessionName: sessionName,
		Windows:     windows,
		Environment: environment,
		filesystem:  fs,
	}
}
