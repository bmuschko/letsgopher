package path

import (
	"os"
	"path/filepath"
)

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
