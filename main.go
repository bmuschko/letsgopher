package main

import (
	"fmt"
	"github.com/bmuschko/letsgopher/cmd"
	"github.com/bmuschko/letsgopher/template/download"
	"os"
)

var version = "undefined"

func main() {
	download.SetVersion(version)
	cmd.SetVersion(version)

	cmd := cmd.NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
