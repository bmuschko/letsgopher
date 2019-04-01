package templ

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"

	"github.com/bmuschko/lets-gopher/utils"
)

var LetsGopherSettings = EnvSettings{Home: Home(filepath.Join(homeDir(), ".letsgopher"))}

type EnvSettings struct {
	Home Home
}

type Home string

func (h Home) String() string {
	return os.ExpandEnv(string(h))
}

func (h Home) ArchiveDir() string {
	return h.path("archive")
}

func (h Home) TemplatesFile() string {
	return h.path("templates.yaml")
}

func (h Home) path(elem ...string) string {
	p := []string{h.String()}
	p = append(p, elem...)
	return filepath.Join(p...)
}

func homeDir() string {
	homeDir, err := homedir.Dir()
	utils.CheckIfError(err)
	return homeDir
}
