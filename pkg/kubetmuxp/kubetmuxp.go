package kubetmuxp

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/thecasualcoder/kube-tmuxp/pkg/filesystem"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubeconfig"
	"github.com/thecasualcoder/kube-tmuxp/pkg/tmuxp"
	yamlV2 "gopkg.in/yaml.v2"
)

// Envs reprensents environemnt variables
type Envs map[string]string

//Cluster represents a Kubernetes cluster
type Cluster struct {
	Name    string `yaml:"name"`
	Zone    string `yaml:"zone"`
	Region  string `yaml:"region"`
	Context string `yaml:"context"`
	Envs    `yaml:"envs"`
}

// DefaultContextName returns default context name
func (c *Cluster) DefaultContextName(project string) (string, error) {
	if regional, err := c.IsRegional(); err != nil {
		return "", err
	} else if regional {
		return fmt.Sprintf("gke_%s_%s_%s", project, c.Region, c.Name), nil
	} else {
		return fmt.Sprintf("gke_%s_%s_%s", project, c.Zone, c.Name), nil
	}
}

// IsRegional tells if a cluster is a regional cluster
func (c *Cluster) IsRegional() (bool, error) {
	if c.Region != "" && c.Zone != "" {
		return false, fmt.Errorf("Only one of region or zone should be given")
	}

	if c.Region != "" && c.Zone == "" {
		return true, nil
	}

	return false, nil
}

// Clusters represents a list of Kubernetes clusters
type Clusters []Cluster

//Project represents a cloud project
type Project struct {
	Name     string `yaml:"name"`
	Clusters `yaml:"clusters"`
}

//Projects represents a list of cloud projects
type Projects []Project

// Config represents kube-tmuxp config
type Config struct {
	Projects   `yaml:"projects"`
	filesystem filesystem.FileSystem
	kubeCfg    kubeconfig.KubeConfig
}

func (c *Config) load(cfgFile string) error {
	reader, err := c.filesystem.Open(cfgFile)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	err = yamlV2.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) saveTmuxpConfig(kubeCfgFile string, cluster Cluster) error {
	windows := tmuxp.Windows{{Name: "default"}}
	env := tmuxp.Environment{"KUBECONFIG": kubeCfgFile}
	for k, v := range cluster.Envs {
		env[k] = v
	}

	tmuxpCfg, err := tmuxp.NewConfig(cluster.Context, windows, env, c.filesystem)
	if err != nil {
		return err
	}

	tmuxpCfgFile := path.Join(tmuxpCfg.TmuxpConfigsDir(), fmt.Sprintf("%s.yaml", cluster.Context))
	if err := tmuxpCfg.Save(tmuxpCfgFile); err != nil {
		return err
	}
	return nil
}

// Process processes kube-tmuxp configs
func (c *Config) Process() error {
	kubeCfgsDir := c.kubeCfg.KubeCfgsDir()
	for _, project := range c.Projects {
		for _, cluster := range project.Clusters {
			kubeCfgFile := path.Join(kubeCfgsDir, cluster.Context)

			fmt.Printf("Cluster: %s\n", cluster.Name)
			fmt.Println("Deleting exisiting context...")
			if err := c.kubeCfg.Delete(kubeCfgFile); err != nil {
				return err
			}

			fmt.Println("Adding context...")
			if regional, err := cluster.IsRegional(); err != nil {
				return err
			} else if regional {
				if err := c.kubeCfg.AddRegionalCluster(project.Name, cluster.Name, cluster.Region, kubeCfgFile); err != nil {
					return err
				}
			} else {
				if err := c.kubeCfg.AddZonalCluster(project.Name, cluster.Name, cluster.Zone, kubeCfgFile); err != nil {
					return err
				}
			}

			fmt.Println("Renaming context...")
			defaultCtxName, err := cluster.DefaultContextName(project.Name)
			if err != nil {
				return err
			}
			c.kubeCfg.RenameContext(defaultCtxName, cluster.Context, kubeCfgFile)

			fmt.Println("Creating tmuxp config...")
			c.saveTmuxpConfig(kubeCfgFile, cluster)

			fmt.Println("")
		}
	}

	return nil
}

// NewConfig creates a new kube-tmuxp Config
func NewConfig(cfgFile string, fs filesystem.FileSystem, kubeCfg kubeconfig.KubeConfig) (Config, error) {
	cfg := Config{
		filesystem: fs,
		kubeCfg:    kubeCfg,
	}

	if err := cfg.load(cfgFile); err != nil {
		return cfg, err
	}

	return cfg, nil
}
