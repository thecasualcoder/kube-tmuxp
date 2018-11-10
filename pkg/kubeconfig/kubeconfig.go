package kubeconfig

import (
	"os"
	"path"

	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
)

// KubeConfig exposes methods to perform actions on kubeconfig
type KubeConfig struct {
	filesystem filesystem.FileSystem
	dir        string
}

// Delete deletes the kubeconfig file for the given context
func (k *KubeConfig) Delete(context string) error {
	file := path.Join(k.dir, context)

	if err := k.filesystem.Remove(file); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// New returns a new KubeConfig
func New(fs filesystem.FileSystem) (KubeConfig, error) {
	home, err := fs.HomeDir()
	if err != nil {
		return KubeConfig{}, err
	}
	kubeConfigsDir := path.Join(home, ".kube/configs")

	return KubeConfig{
		filesystem: fs,
		dir:        kubeConfigsDir,
	}, nil
}
