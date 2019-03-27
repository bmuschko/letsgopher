package main

import (
	"github.com/bmuschko/lets-gopher/cmd"
	"github.com/bmuschko/lets-gopher/templ/download"
)

var (
	version = "undefined"
)

func main() {
	download.SetVersion(version)
	cmd.SetVersion(version)
	cmd.Execute()
}
