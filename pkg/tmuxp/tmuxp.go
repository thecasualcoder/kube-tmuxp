package tmuxp

import "io"

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
	writer io.Writer
}

// New returns a new tmuxp config
func New(writer io.Writer) Config {
	return Config{
		writer: writer,
	}
}
