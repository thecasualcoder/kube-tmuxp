package kubeconfig

import (
	"fmt"
	"os"
	"path"

	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
	homedir "github.com/mitchellh/go-homedir"
)

// KubeConfig exposes methods to perform actions on kubeconfig
type KubeConfig struct {
	filesystem filesystem.FileSystem
	dir        string
}

// Delete deletes the kubeconfig file for the given context
func (k *KubeConfig) Delete(context string) error {
	file := path.Join(k.dir, context)

	if err := k.filesystem.Remove(file); err != nil {
		return err
	}

	return nil
}

// New returns a new KubeConfig
func New(fs filesystem.FileSystem) KubeConfig {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	kubeConfigsDir := path.Join(home, ".kube/configs")

	return KubeConfig{
		filesystem: fs,
		dir:        kubeConfigsDir,
	}
}
