package main

import (
	"github.com/nerveband/mochi-cli/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
