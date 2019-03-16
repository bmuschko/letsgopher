package templ

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"time"
)

type TemplatesFile struct {
	Generated time.Time   `json:"generated"`
	Templates []*Template `json:"templates"`
}

type Template struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	ArchivePath string `json:"archivePath"`
}

func NewTemplatesFile() *TemplatesFile {
	return &TemplatesFile{
		Generated: time.Now(),
		Templates: []*Template{},
	}
}

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

func (r *TemplatesFile) Has(name string, version string) bool {
	for _, rf := range r.Templates {
		if rf.Name == name && rf.Version == version {
			return true
		}
	}
	return false
}

func (r *TemplatesFile) Get(name string, version string) *Template {
	for _, rf := range r.Templates {
		if rf.Name == name && rf.Version == version {
			return rf
		}
	}
	return nil
}

func (r *TemplatesFile) Add(re ...*Template) {
	r.Templates = append(r.Templates, re...)
}

func (r *TemplatesFile) Update(re ...*Template) {
	for _, target := range re {
		found := false
		for j, template := range r.Templates {
			if template.Name == target.Name {
				r.Templates[j] = target
				found = true
				break
			}
		}
		if !found {
			r.Add(target)
		}
	}
}

func (r *TemplatesFile) Remove(name string) bool {
	cp := []*Template{}
	found := false
	for _, rf := range r.Templates {
		if rf.Name == name {
			found = true
			continue
		}
		cp = append(cp, rf)
	}
	r.Templates = cp
	return found
}

func (r *TemplatesFile) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}
