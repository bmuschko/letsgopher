package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"time"
)

// TemplatesFile represents a local template registry file.
type TemplatesFile struct {
	Generated time.Time   `json:"generated"`
	Templates []*Template `json:"templates"`
}

// Template represents a template in the local registry.
type Template struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	ArchivePath string `json:"archivePath"`
}

// NewTemplatesFile creates a local template registry file of type TemplatesFile.
func NewTemplatesFile() *TemplatesFile {
	return &TemplatesFile{
		Generated: time.Now(),
		Templates: []*Template{},
	}
}

// LoadTemplatesFile loads the template registry file.
func LoadTemplatesFile(path string) (*TemplatesFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := &TemplatesFile{}
	err = yaml.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Has checks if a template with a given name and version has been installed.
func (r *TemplatesFile) Has(name string, version string) bool {
	for _, rf := range r.Templates {
		if rf.Name == name && rf.Version == version {
			return true
		}
	}
	return false
}

// Get retrieves a template with a given name and version from the registry.
func (r *TemplatesFile) Get(name string, version string) *Template {
	for _, rf := range r.Templates {
		if rf.Name == name && rf.Version == version {
			return rf
		}
	}
	return nil
}

// Add adds a template to the registry.
func (r *TemplatesFile) Add(re ...*Template) {
	r.Templates = append(r.Templates, re...)
}

// Update updates an existing template in the registry.
func (r *TemplatesFile) Update(re ...*Template) bool {
	found := false
	for _, target := range re {
		for j, template := range r.Templates {
			if template.Name == target.Name && template.Version == target.Version {
				r.Templates[j] = target
				found = true
				break
			}
		}
		if !found {
			r.Add(target)
		}
	}
	return found
}

// Remove removes an existing template from the registry.
func (r *TemplatesFile) Remove(name string, version string) bool {
	cp := []*Template{}
	found := false
	for _, rf := range r.Templates {
		if rf.Name == name && rf.Version == version {
			found = true
			continue
		}
		cp = append(cp, rf)
	}
	r.Templates = cp
	return found
}

// WriteFile write the template registry file.
func (r *TemplatesFile) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}
