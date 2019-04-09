package cmd

import (
	"bytes"
	"fmt"
	"github.com/Flaque/filet"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestInspectNonExistentTemplateFile(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	b := bytes.NewBuffer(nil)
	templateInspect := &templateInspectCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	err := templateInspect.run()

	assert.NotNil(t, err)
	assert.Equal(t, "failed to load templates.yaml file", err.Error())
}

func TestInspectNonExistentTemplateName(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	b := bytes.NewBuffer(nil)
	templateInspect := &templateInspectCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		out:             b,
		home:            storage.Home(tmpHome),
	}
	templatesFile := storage.Home(tmpHome).TemplatesFile()
	f, err := os.Create(templatesFile)
	if err != nil {
		t.Errorf("failed to create file %s", templatesFile)
	}
	_, err = f.WriteString(`generated: "2019-03-15T16:31:57.232715-06:00"
templates: []`)
	if err != nil {
		t.Errorf("failed to write to file %s", f.Name())
	}
	defer f.Close()
	err = templateInspect.run()

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("template with name %q and version %q hasn't been installed", "hello-world", "1.0.0"), err.Error())
}

func TestInspectValidTemplate(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	b := bytes.NewBuffer(nil)
	aM := new(ArchiverMock)
	templateInspect := &templateInspectCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		out:             b,
		home:            storage.Home(tmpHome),
		archiver:        aM,
	}
	templatesFile := storage.Home(tmpHome).TemplatesFile()
	f, err := os.Create(templatesFile)
	if err != nil {
		t.Errorf("failed to create file %s", templatesFile)
	}
	archiveFile := fmt.Sprintf("%s/archive/hello-world-1.0.0.zip", tmpHome)
	_, err = f.WriteString(fmt.Sprintf(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: %s
  name: hello-world
  version: 1.0.0`, archiveFile))
	if err != nil {
		t.Errorf("failed to write to file %s", f.Name())
	}
	defer f.Close()
	aM.On("LoadManifestFile", archiveFile).Return([]byte(`version: "1.0.0"
parameters:
	- name: "module"
prompt: "Please provide a module name"
	type: "string"
description: "The module name is used in the go.mod file"
	- name: "message"
prompt: "Please select a message"
	type: "string"
enum: ["Hello World!", "Let's get started", "This is just the beginning"]
	description: "The message to be rendered when executing the program"`), nil)
	err = templateInspect.run()

	aM.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, `template:
  name: "hello-world"
  version: "1.0.0"
manifest:
  version: "1.0.0"
  parameters:
  	- name: "module"
  prompt: "Please provide a module name"
  	type: "string"
  description: "The module name is used in the go.mod file"
  	- name: "message"
  prompt: "Please select a message"
  	type: "string"
  enum: ["Hello World!", "Let's get started", "This is just the beginning"]
  	description: "The message to be rendered when executing the program"
`, b.String())
}

type ArchiverMock struct {
	mock.Mock
}

func (a *ArchiverMock) Extract(archiveFile string, targetDir string, replacements map[string]interface{}) error {
	args := a.Called(archiveFile, targetDir, replacements)
	return args.Error(0)
}

func (a *ArchiverMock) LoadManifestFile(src string) ([]byte, error) {
	args := a.Called(src)
	return args.Get(0).([]byte), args.Error(1)
}
