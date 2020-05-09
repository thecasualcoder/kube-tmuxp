package generator

import (
	"fmt"
	"github.com/thecasualcoder/kube-tmuxp/pkg/commander"
	"github.com/thecasualcoder/kube-tmuxp/pkg/file"
	"github.com/thecasualcoder/kube-tmuxp/pkg/filesystem"
	"github.com/thecasualcoder/kube-tmuxp/pkg/gcloud"
	"io"
)

type Generator interface {
	Generate(outStream, errStream io.Writer)
}

type Options struct {
	From           string
	AllProjects    bool
	ProjectIDs     []string
	AdditionalEnvs []string
	Apply          bool
	CfgFile        string
}

func NewGenerator(options Options, fs filesystem.FileSystem, cmdr commander.Commander) (Generator, error) {
	switch options.From {
	case "file":
		if err := areFlagsValidForSourceFile(options.AllProjects, options.ProjectIDs, options.AdditionalEnvs); err != nil {
			return nil, fmt.Errorf("error in the flags for source type 'file': %s", err)
		}
		return file.NewGenerator(fs, cmdr, options.CfgFile), nil
	case "gcloud":
		return gcloud.NewGenerator(options.ProjectIDs, options.AllProjects, options.AdditionalEnvs, options.Apply), nil
	default:
		return nil, fmt.Errorf("invalid source provided: valid sources are file,gcloud")
	}
}

func areFlagsValidForSourceFile(allProjects bool, projectIDs, additionalEnvs []string) error {
	err := ""
	counter := 1
	if projectIDs != nil {
		err += fmt.Sprintf("\n %d) %s", counter, "project-ids should be empty for source file")
		counter++
	}
	if allProjects {
		err += fmt.Sprintf("\n %d) %s", counter, "all-projects should be false for source file")
		counter++
	}
	if additionalEnvs != nil {
		err += fmt.Sprintf("\n %d) %s", counter, "additional-envs should be empty for source file")
	}

	if err != "" {
		return fmt.Errorf("%s\n", err)
	}
	return nil
}
