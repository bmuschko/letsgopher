package storage

import (
	"os"
	"path/filepath"
)

// This helper builds paths relative to a Home directory.
type Home string

// String returns Home as a string.
//
// Implements fmt.Stringer.
func (h Home) String() string {
	return os.ExpandEnv(string(h))
}

// ArchiveDir returns the path to the archive directory.
func (h Home) ArchiveDir() string {
	return h.path("archive")
}

// TemplatesFile returns the path to the templates registry file.
func (h Home) TemplatesFile() string {
	return h.path("templates.yaml")
}

func (h Home) path(elem ...string) string {
	p := []string{h.String()}
	p = append(p, elem...)
	return filepath.Join(p...)
}
