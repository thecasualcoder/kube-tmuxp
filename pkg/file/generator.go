package file

import (
	"fmt"
	"io"
	"os"

	"github.com/thecasualcoder/kube-tmuxp/pkg/commander"
	"github.com/thecasualcoder/kube-tmuxp/pkg/filesystem"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubeconfig"
	"github.com/thecasualcoder/kube-tmuxp/pkg/kubetmuxp"
)

type Generator struct {
	fs      filesystem.FileSystem
	cmdr    commander.Commander
	cfgFile string
}

func NewGenerator(fs filesystem.FileSystem, cmdr commander.Commander, cfgFile string) Generator {
	return Generator{fs: fs, cmdr: cmdr, cfgFile: cfgFile}
}

func (g Generator) Generate(_, errStream io.Writer) {
	kubeCfg, err := kubeconfig.New(g.fs, g.cmdr)
	if err != nil {
		_, _ = fmt.Fprintf(errStream, err.Error())
		os.Exit(1)
	}

	fmt.Println("Using config file:", g.cfgFile)
	kubetmuxpCfg, err := kubetmuxp.NewConfig(g.cfgFile, g.fs, kubeCfg)
	if err != nil {
		_, _ = fmt.Fprintf(errStream, err.Error())
		os.Exit(1)
	}

	if err = kubetmuxpCfg.Process(); err != nil {
		_, _ = fmt.Fprintf(errStream, err.Error())
		os.Exit(1)
	}
}
