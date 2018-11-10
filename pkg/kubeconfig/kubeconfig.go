package kubeconfig

import (
	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
)

// KubeConfig exposes methods to perform actions on kubeconfig
type KubeConfig struct {
	filesystem filesystem.FileSystem
}

// Delete deletes the given kubeconfig file
func (k *KubeConfig) Delete(file string) error {
	if err := k.filesystem.Remove(file); err != nil {
		return err
	}

	return nil
}

// New returns a new KubeConfig
func New(fs filesystem.FileSystem) KubeConfig {
	return KubeConfig{
		filesystem: fs,
	}
}
