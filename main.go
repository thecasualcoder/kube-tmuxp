package main

import "github.com/jfreeland/kube-tmuxp/cmd"

var version string

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
