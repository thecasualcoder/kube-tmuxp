package filesystem

import (
	"bufio"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

// FileSystem represents a filesystem
type FileSystem interface {
	Remove(file string) error
	HomeDir() (string, error)
	Open(file string) (io.Reader, error)
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

// Open opens a file for reading
func (d *Default) Open(file string) (io.Reader, error) {
	reader, err := os.Open(file)
	if err != nil {
		return reader, err
	}
	return bufio.NewReader(reader), err
}
