package templ

import (
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/ghodss/yaml"
	"strconv"
)

const (
	maxCompatManifestVersion = "1.0.0"
	StringType               = "string"
	IntegerType              = "integer"
	BooleanType              = "boolean"
)

type ManifestFile struct {
	Version    string       `json:"version"`
	Parameters []*Parameter `json:"parameters"`
}

type Parameter struct {
	Name         string   `json:"name"`
	Prompt       string   `json:"prompt"`
	Type         string   `json:"type"`
	Enum         []string `json:"enum"`
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
}

func LoadManifestData(b []byte) (*ManifestFile, error) {
	m := &ManifestFile{}
	err := yaml.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}

	err = validateManifest(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func validateManifest(m *ManifestFile) error {
	err := validateManifestVersion(m.Version)
	if err != nil {
		return err
	}
	err = validateManifestParams(m.Parameters)
	if err != nil {
		return err
	}
	return nil
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

func validateManifestParams(params []*Parameter) error {
	for _, p := range params {
		if p.Type == IntegerType {
			_, err := strconv.Atoi(p.DefaultValue)
			if err != nil {
				return err
			}
			if p.Enum != nil {
				for _, e := range p.Enum {
					_, err := strconv.Atoi(e)
					if err != nil {
						return err
					}
				}
			}
		}
		if p.Type == BooleanType {
			_, err := strconv.ParseBool(p.DefaultValue)
			if err != nil {
				return err
			}
			if p.Enum != nil {
				return errors.New("boolean type does not allow enums")
			}
		}
	}
	return nil
}
