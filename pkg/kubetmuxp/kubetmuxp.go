package kubetmuxp

import (
	"io"
	"io/ioutil"

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
func (c *Config) Process() error {
	for _, project := range c.Projects {
		for _, cluster := range project.Clusters {
			if err := c.kubeCfg.Delete(cluster.Context); err != nil {
				return err
			}
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
