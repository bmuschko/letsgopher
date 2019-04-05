package environment

import (
	"github.com/bmuschko/lets-gopher/template/storage"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

var defaultHome = filepath.Join(homeDir(), ".letsgopher")

// Settings exposes the environment settings.
var Settings = EnvSettings{Home: storage.Home(defaultHome)}

// EnvSettings describes the environment settings.
type EnvSettings struct {
	Home storage.Home
}

func homeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
