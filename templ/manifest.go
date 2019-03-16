package templ

import (
	"github.com/ghodss/yaml"
)

type ManifestFile struct {
	Parameters []*Parameter `json:"parameters"`
}

type Parameter struct {
	Name         string `json:"name"`
	Prompt       string `json:"prompt"`
	Type         string `json:"type"`
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
}

func LoadManifestData(b []byte) (*ManifestFile, error) {
	m := &ManifestFile{}
	err := yaml.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
