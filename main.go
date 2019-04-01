package main

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/cmd"
	"github.com/bmuschko/lets-gopher/templ/download"
	"os"
)

var (
	version = "undefined"
)

func main() {
	download.SetVersion(version)
	cmd.SetVersion(version)

	cmd := cmd.NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
