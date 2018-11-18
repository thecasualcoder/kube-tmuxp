package tmuxp

import (
	"path"

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
	filesystem   filesystem.FileSystem
	tmuxpCfgsDir string
}

// TmuxpConfigsDir returns the directory in which tmuxp
// configs are stored
func (c *Config) TmuxpConfigsDir() string {
	return c.tmuxpCfgsDir
}

// NewConfig returns a new tmuxp config
func NewConfig(sessionName string, windows Windows, environment Environment, fs filesystem.FileSystem) (*Config, error) {
	home, err := fs.HomeDir()
	if err != nil {
		return nil, err
	}
	tmuxpCfgsDir := path.Join(home, ".tmuxp")

	return &Config{
		SessionName:  sessionName,
		Windows:      windows,
		Environment:  environment,
		filesystem:   fs,
		tmuxpCfgsDir: tmuxpCfgsDir,
	}, nil
}
