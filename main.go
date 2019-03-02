package main

import (
	"github.com/bmuschko/lets-gopher/cmd"
	"github.com/bmuschko/lets-gopher/templ"
)

var (
	version = "undefined"
)

func main() {
	templ.SetVersion(version)
	cmd.SetVersion(version)
	cmd.Execute()
}
