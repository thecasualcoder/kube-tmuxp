package kubetmuxp

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Envs abcd
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

// NewConfig creates a new Config
func NewConfig(reader io.Reader) Config {
	return Config{
		reader: reader,
	}
}
