package main

import "github.com/thecasualcoder/kube-tmuxp/cmd"

var version string

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
