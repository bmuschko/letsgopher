package environment

import (
	"github.com/bmuschko/lets-gopher/templ/path"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

var defaultHome = filepath.Join(homeDir(), ".letsgopher")
var Settings = EnvSettings{Home: path.Home(defaultHome)}

type EnvSettings struct {
	Home path.Home
}

func homeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
