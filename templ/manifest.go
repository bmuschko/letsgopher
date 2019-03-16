package templ

import (
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/ghodss/yaml"
)

const maxCompatManifestVersion = "1.0.0"

type ManifestFile struct {
	Version    string       `json:"version"`
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

	err = validateManifestVersion(m.Version)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func validateManifestVersion(version string) error {
	if version == "" {
		return errors.New("manifest file needs to provide a version")
	}

	v, err := semver.Make(version)
	if err != nil {
		return err
	}

	maxCompatVersion, err := semver.Make(maxCompatManifestVersion)
	if err != nil {
		return err
	}

	if v.GT(maxCompatVersion) {
		return fmt.Errorf("manifest version needs to be less than %s", maxCompatManifestVersion)
	}

	return nil
}
