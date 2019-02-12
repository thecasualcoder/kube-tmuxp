package tmuxp

import (
	"path"

	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
	yaml "gopkg.in/yaml.v2"
)

// Window represents a window in a tmux session
type Window struct {
	Name  string     `yaml:"window_name"`
	Panes []struct{} `yaml:"panes"`
}

// Windows represents a list of windows in a tmux session
type Windows []Window

// Environment represents env variables to be loaded in a tmux session
type Environment map[string]string

// Config represents tmuxp config
type Config struct {
	SessionName  string `yaml:"session_name"`
	Windows      `yaml:"windows"`
	Environment  `yaml:"environment"`
	filesystem   filesystem.FileSystem
	tmuxpCfgsDir string
}

// TmuxpConfigsDir returns the directory in which tmuxp
// configs are stored
func (c *Config) TmuxpConfigsDir() string {
	return c.tmuxpCfgsDir
}

// Save saves the tmuxp config as file
func (c *Config) Save(file string) error {
	writer, err := c.filesystem.Create(file)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

// NewConfig returns a new tmuxp config
func NewConfig(sessionName string, windows Windows, environment Environment, fs filesystem.FileSystem) (*Config, error) {
	home, err := fs.HomeDir()
	if err != nil {
		return nil, err
	}
	tmuxpCfgsDir := path.Join(home, ".tmuxp")

	err = fs.CreateDirIfNotExist(tmuxpCfgsDir)
	if err != nil {
		return nil, err
	}
	return &Config{
		SessionName:  sessionName,
		Windows:      windows,
		Environment:  environment,
		filesystem:   fs,
		tmuxpCfgsDir: tmuxpCfgsDir,
	}, nil
}
