package manifest

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadManifestDataForCorrectDefinition(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	content := []byte(`version: "0.1.0"
parameters:
  - name: "module"
    prompt: "Please provide a module name"
    type: "string"
    description: "The module name is used in the go.mod file"
  - name: "message"
    prompt: "Please select a message"
    type: "string"
    enum: ["Hello World!", "Let's get started", "This is just the beginning"]
    description: "The message to be rendered when executing the program"`)
	manifestFile, err := LoadManifestData(content)
	param1 := Parameter{
		Name:        "module",
		Prompt:      "Please provide a module name",
		Type:        "string",
		Description: "The module name is used in the go.mod file",
	}
	param2 := Parameter{
		Name:        "message",
		Prompt:      "Please select a message",
		Type:        "string",
		Enum:        []string{"Hello World!", "Let's get started", "This is just the beginning"},
		Description: "The message to be rendered when executing the program",
	}

	assert.NotNil(t, manifestFile)
	assert.Nil(t, err)
	assert.Equal(t, "0.1.0", manifestFile.Version)
	assert.Equal(t, 2, len(manifestFile.Parameters))
	assert.Equal(t, param1, *manifestFile.Parameters[0])
	assert.Equal(t, param2, *manifestFile.Parameters[1])
}

func TestLoadManifestDataForIncorrectDefinition(t *testing.T) {
	content := []byte("test")
	manifestFile, err := LoadManifestData(content)

	assert.Nil(t, manifestFile)
	assert.NotNil(t, err)
	assert.Equal(t, "error unmarshaling JSON: json: cannot unmarshal string into Go value of type manifest.ManifestFile", err.Error())
}

func TestValidateManifestWithEmptyVersion(t *testing.T) {
	manifestFile := &ManifestFile{Version: ""}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "manifest file needs to provide a version", err.Error())
}

func TestValidateManifestWithIncorrectSemVerVersion(t *testing.T) {
	semVers := []invalidSemVer{
		{"a.b.c", "Invalid character(s) found in major number \"a\""},
		{"1", "No Major.Minor.Patch elements found"},
		{"1.2", "No Major.Minor.Patch elements found"},
		{"1.2.", "strconv.ParseUint: parsing \"\": invalid syntax"},
	}
	for _, sv := range semVers {
		t.Run(sv.version, func(t *testing.T) {
			manifestFile := &ManifestFile{Version: sv.version}
			err := ValidateManifest(manifestFile)

			assert.NotNil(t, err)
			assert.Equal(t, sv.errorMessage, err.Error())
		})
	}
}

func TestValidateManifestWithVersionOutOfSupportedRange(t *testing.T) {
	manifestFile := &ManifestFile{Version: "2.4.5"}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "manifest version needs to be less than 1.0.0", err.Error())
}

func TestValidateManifestWithMaxSupportedVersion(t *testing.T) {
	manifestFile := &ManifestFile{Version: "1.0.0"}
	err := ValidateManifest(manifestFile)

	assert.Nil(t, err)
}

func TestValidateManifestWithParameters(t *testing.T) {
	manifestFile := &ManifestFile{Version: "1.0.0"}
	err := ValidateManifest(manifestFile)

	assert.Nil(t, err)
}

func TestValidateManifestWithEmptyParameterType(t *testing.T) {
	manifestFile := &ManifestFile{
		Version: "1.0.0",
		Parameters: []*Parameter{
			{Type: ""},
		},
	}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "every parameter defined in manifest needs to provide a type", err.Error())
}

func TestValidateManifestWithIncorrectIntegerParameterDefaultValue(t *testing.T) {
	manifestFile := &ManifestFile{
		Version: "1.0.0",
		Parameters: []*Parameter{
			{Type: "integer", DefaultValue: "abc"},
		},
	}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "strconv.Atoi: parsing \"abc\": invalid syntax", err.Error())
}

func TestValidateManifestWithIncorrectIntegerEnumValues(t *testing.T) {
	manifestFile := &ManifestFile{
		Version: "1.0.0",
		Parameters: []*Parameter{
			{Type: "integer", Enum: []string{"123", "abc"}},
		},
	}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "strconv.Atoi: parsing \"abc\": invalid syntax", err.Error())
}

func TestValidateManifestWithIncorrectBooleanParameterDefaultValue(t *testing.T) {
	manifestFile := &ManifestFile{
		Version: "1.0.0",
		Parameters: []*Parameter{
			{Type: "boolean", DefaultValue: "notaboolean"},
		},
	}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "strconv.ParseBool: parsing \"notaboolean\": invalid syntax", err.Error())
}

func TestValidateManifestWithBooleanEnumValues(t *testing.T) {
	manifestFile := &ManifestFile{
		Version: "1.0.0",
		Parameters: []*Parameter{
			{Type: "boolean", Enum: []string{"true", "false"}},
		},
	}
	err := ValidateManifest(manifestFile)

	assert.NotNil(t, err)
	assert.Equal(t, "boolean type does not allow enums", err.Error())
}

type invalidSemVer struct {
	version      string
	errorMessage string
}
