package main

import "github.com/arunvelsriram/kube-tmuxp/cmd"

var version string

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
