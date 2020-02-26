package kubeconfig

import (
	"fmt"
	"os"
	"path"

	"github.com/jfreeland/kube-tmuxp/pkg/commander"
	"github.com/jfreeland/kube-tmuxp/pkg/filesystem"
)

// KubeConfig exposes methods to perform actions on kubeconfig
type KubeConfig struct {
	filesystem  filesystem.FileSystem
	commander   commander.Commander
	kubeCfgsDir string
}

// Delete deletes the given kubeconfig file
func (k *KubeConfig) Delete(kubeCfgFile string) error {
	if err := k.filesystem.Remove(kubeCfgFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// AddRegionalEKSCluster imports Kubernetes context for a regional EKS cluster
func (k *KubeConfig) AddRegionalEKSCluster(cluster string, region string, role string, env map[string]string, kubeCfgFile string) error {
	var args []string
	if len(role) > 0 {
		args = []string{
			"eks",
			"update-kubeconfig",
			"--name",
			cluster,
			"--kubeconfig",
			kubeCfgFile,
			"--region",
			region,
			"--role-arn",
			role,
		}
	} else {
		args = []string{
			"eks",
			"update-kubeconfig",
			"--name",
			cluster,
			"--kubeconfig",
			kubeCfgFile,
			"--region",
			region,
		}
	}

	envs := []string{}
	for k, v := range env {
		envs = append(envs, k+"="+v)
	}
	if _, err := k.commander.Execute("aws", args, envs); err != nil {
		return err
	}

	return nil
}

// AddRegionalGKECluster imports Kubernetes context for
// a regional Kubernetes cluster
func (k *KubeConfig) AddRegionalGKECluster(project string, cluster string, region string, kubeCfgFile string) error {
	args := []string{
		"beta",
		// Is beta needed here?  You can run into an issue where if
		// you haven't updated gcloud lately, beta will stop and
		// report that it needs updated, but gcloud itself will not.
		"container",
		"clusters",
		"get-credentials",
		cluster,
		fmt.Sprintf("--region=%s", region),
		fmt.Sprintf("--project=%s", project),
	}
	envs := []string{
		fmt.Sprintf("KUBECONFIG=%s", kubeCfgFile),
	}
	if _, err := k.commander.Execute("gcloud", args, envs); err != nil {
		return err
	}

	return nil
}

// AddZonalGKECluster imports Kubernetes context for
// a zonal Kubernetes cluster
func (k *KubeConfig) AddZonalGKECluster(project string, cluster string, zone string, kubeCfgFile string) error {
	args := []string{
		"container",
		"clusters",
		"get-credentials",
		cluster,
		fmt.Sprintf("--zone=%s", zone),
		fmt.Sprintf("--project=%s", project),
	}
	envs := []string{
		fmt.Sprintf("KUBECONFIG=%s", kubeCfgFile),
	}
	if _, err := k.commander.Execute("gcloud", args, envs); err != nil {
		return err
	}

	return nil
}

// RenameContext renames a Kubernetes context
func (k *KubeConfig) RenameContext(oldCtx string, newCtx string, kubeCfgFile string) error {
	args := []string{
		"config",
		"rename-context",
		oldCtx,
		newCtx,
	}
	envs := []string{
		fmt.Sprintf("KUBECONFIG=%s", kubeCfgFile),
	}
	if _, err := k.commander.Execute("kubectl", args, envs); err != nil {
		return err
	}

	return nil
}

// KubeCfgsDir returns the directory in which kube configs are stored
func (k *KubeConfig) KubeCfgsDir() string {
	return k.kubeCfgsDir
}

// New returns a new KubeConfig
func New(fs filesystem.FileSystem, cmdr commander.Commander) (KubeConfig, error) {
	home, err := fs.HomeDir()
	if err != nil {
		return KubeConfig{}, err
	}
	kubeConfigsDir := path.Join(home, ".kube/configs")

	return KubeConfig{
		filesystem:  fs,
		commander:   cmdr,
		kubeCfgsDir: kubeConfigsDir,
	}, nil
}
