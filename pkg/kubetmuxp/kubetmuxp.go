package kubetmuxp

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config represents kube-tmuxp config
type Config struct {
	Projects []struct {
		Name     string `yaml:"name"`
		Clusters []struct {
			Name    string            `yaml:"name"`
			Zone    string            `yaml:"zone"`
			Region  string            `yaml:"region"`
			Context string            `yaml:"context"`
			Envs    map[string]string `yaml:"envs"`
		}
	}
}

// NewConfig constructs kube-tmuxp config from given file
func NewConfig(cfgFile string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
