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

// Role represents a Role ARN required to connect to an EKS cluster
type Role string

// ProviderEnvs represents environment variables for the provider
type ProviderEnvs map[string]string

// Envs reprensents environemnt variables for Tmux
type TmuxEnvs map[string]string

//Cluster represents a Kubernetes cluster
type Cluster struct {
	Name         string `yaml:"name"`
	Zone         string `yaml:"zone"`
	Region       string `yaml:"region"`
	Context      string `yaml:"context"`
	Role         string `yaml:"role"`
	TmuxEnvs     `yaml:"tmux_envs"`
	ProviderEnvs `yaml:"provider_envs"`
}

// DefaultEKSContextName returns default context name for EKS
func (c *Cluster) DefaultEKSContextName(project string) (string, error) {
	return fmt.Sprintf("arn:aws:eks:%s:%s:cluster/%s", c.Region, project, c.Name), nil
}

// DefaultGKEContextName returns default context name for GKE
func (c *Cluster) DefaultGKEContextName(project string) (string, error) {
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
	Provider string `yaml:"provider"`
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
	for k, v := range cluster.TmuxEnvs {
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

			if project.Provider == "" {
				project.Provider = "gke"
			}
			fmt.Printf("Provider: %s\n", project.Provider)
			fmt.Printf("Cluster: %s\n", cluster.Name)
			fmt.Println("Deleting exisiting context...")
			if err := c.kubeCfg.Delete(kubeCfgFile); err != nil {
				return err
			}

			fmt.Println("Adding context...")
			if project.Provider == "gke" {
				if regional, err := cluster.IsRegional(); err != nil {
					return err
				} else if regional {
					if err := c.kubeCfg.AddRegionalGKECluster(project.Name, cluster.Name, cluster.Region, kubeCfgFile); err != nil {
						return err
					}
				} else {
					if err := c.kubeCfg.AddZonalGKECluster(project.Name, cluster.Name, cluster.Zone, kubeCfgFile); err != nil {
						return err
					}
				}
			} else if project.Provider == "eks" {
				if regional, err := cluster.IsRegional(); err != nil {
					return err
				} else if regional {
					if err := c.kubeCfg.AddRegionalEKSCluster(cluster.Name, cluster.Region, cluster.Role, cluster.ProviderEnvs, kubeCfgFile); err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("Provider must be gke or eks")
			}

			fmt.Println("Renaming context...")
			if project.Provider == "gke" {
				defaultCtxName, err := cluster.DefaultGKEContextName(project.Name)
				if err != nil {
					return err
				}
				c.kubeCfg.RenameContext(defaultCtxName, cluster.Context, kubeCfgFile)
			} else if project.Provider == "eks" {
				defaultCtxName, err := cluster.DefaultEKSContextName(project.Name)
				if err != nil {
					return err
				}
				c.kubeCfg.RenameContext(defaultCtxName, cluster.Context, kubeCfgFile)
			}

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
