package kubetmuxp

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/arunvelsriram/kube-tmuxp/pkg/kubeconfig"
	yaml "gopkg.in/yaml.v2"
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
func (c Cluster) DefaultContextName(project string) (string, error) {
	if regional, err := c.IsRegional(); err != nil {
		return "", err
	} else if regional {
		return fmt.Sprintf("gke_%s_%s_%s", project, c.Region, c.Name), nil
	} else {
		return fmt.Sprintf("gke_%s_%s_%s", project, c.Zone, c.Name), nil
	}
}

// IsRegional tells if a cluster is a regional cluster
func (c Cluster) IsRegional() (bool, error) {
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
	Projects `yaml:"projects"`
	reader   io.Reader
	kubeCfg  kubeconfig.KubeConfig
}

// Load constructs kube-tmuxp config from given file
func (c *Config) Load() error {
	data, err := ioutil.ReadAll(c.reader)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

// Process processes kube-tmuxp configs
func (c Config) Process() error {
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
				c.kubeCfg.AddRegionalCluster(project.Name, cluster.Name, cluster.Region, kubeCfgFile)
			} else {
				c.kubeCfg.AddZonalCluster(project.Name, cluster.Name, cluster.Zone, kubeCfgFile)
			}

			fmt.Println("Renaming context...")

			fmt.Println("")
		}
	}

	return nil
}

// New creates a new kube-tmuxp Config
func New(reader io.Reader, kubeCfg kubeconfig.KubeConfig) Config {
	return Config{
		reader:  reader,
		kubeCfg: kubeCfg,
	}
}
