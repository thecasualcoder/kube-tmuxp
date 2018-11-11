package kubeconfig

import (
	"fmt"
	"os"
	"path"

	"github.com/arunvelsriram/kube-tmuxp/pkg/commander"
	"github.com/arunvelsriram/kube-tmuxp/pkg/filesystem"
)

// KubeConfig exposes methods to perform actions on kubeconfig
type KubeConfig struct {
	filesystem filesystem.FileSystem
	commander  commander.Commander
	dir        string
}

// Delete deletes the kubeconfig file for the given context
func (k KubeConfig) Delete(context string) error {
	file := path.Join(k.dir, context)

	if err := k.filesystem.Remove(file); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// AddRegionalCluster imports Kubernetes context for
// regional Kubernetes cluster
func (k KubeConfig) AddRegionalCluster(project string, cluster string, region string, context string) error {
	kubeconfig := path.Join(k.dir, context)
	envs := []string{
		"CLOUDSDK_CONTAINER_USE_V1_API_CLIENT=false",
		"CLOUDSDK_CONTAINER_USE_V1_API=false",
		fmt.Sprintf("KUBECONFIG=%s", kubeconfig),
	}
	args := []string{
		"beta",
		"container",
		"clusters",
		"get-credentials",
		cluster,
		fmt.Sprintf("--region=%s", region),
		fmt.Sprintf("--project=%s", project),
	}
	out, err := k.commander.Execute("gcloud", args, envs)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}

// AddZonalCluster imports Kubernetes context for
// zonal Kubernetes cluster
func AddZonalCluster(project string, cluster string, zone string, context string) error {
	return nil
}

// New returns a new KubeConfig
func New(fs filesystem.FileSystem, cmdr commander.Commander) (KubeConfig, error) {
	home, err := fs.HomeDir()
	if err != nil {
		return KubeConfig{}, err
	}
	kubeConfigsDir := path.Join(home, ".kube/configs")

	return KubeConfig{
		filesystem: fs,
		commander:  cmdr,
		dir:        kubeConfigsDir,
	}, nil
}
