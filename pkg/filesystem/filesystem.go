package filesystem

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

// FileSystem represents a filesystem
type FileSystem interface {
	Remove(file string) error
	HomeDir() (string, error)
}

// Default represents the Operating System's filesystem
type Default struct{}

// Remove removes a file from Default filesystem
func (d *Default) Remove(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}

// HomeDir returns the home directory
func (d *Default) HomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return home, err
	}

	return home, nil
}
