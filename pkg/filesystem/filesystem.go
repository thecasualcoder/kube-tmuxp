package filesystem

import "os"

// FileSystem abcd
type FileSystem interface {
	Remove(file string) error
}

type defaultFS struct{}

// Remove abcd
func (d *defaultFS) Remove(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}
