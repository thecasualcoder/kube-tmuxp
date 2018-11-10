package filesystem

import "os"

// FileSystem represents a filesystem
type FileSystem interface {
	Remove(file string) error
}

// Default represents the Operating System's filesystem
type Default struct{}

// Remove removes a file from Default filesystem
func (Default) Remove(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}
