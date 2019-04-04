package cmd

import (
	"bytes"
	"github.com/bmuschko/lets-gopher/templ/path"
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCreateProjectWithoutRegisteredTemplate(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	f := path.Home(tmpHome).TemplatesFile()
	err := ioutil.WriteFile(f, []byte(`generated: "2019-03-21T08:49:27.10175-06:00"
templates: []`), 0644)
	if err != nil {
		t.Error("could not write template file")
	}

	b := bytes.NewBuffer(nil)
	projectCreate := &projectCreateCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		out:             b,
		home:            path.Home(tmpHome),
	}
	err = projectCreate.run()

	assert.NotNil(t, err)
	assert.Equal(t, "template with name \"hello-world\" and version \"1.0.0\" hasn't been installed", err.Error())
}

func TestCreateProjectWithRegisteredTemplate(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	f := path.Home(tmpHome).TemplatesFile()
	err := ioutil.WriteFile(f, []byte(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: /my/path/new-project/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0`), 0644)
	if err != nil {
		t.Error("could not write template file")
	}

	b := bytes.NewBuffer(nil)
	aM := new(ArchiverMock)
	projectCreate := &projectCreateCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		targetDir:       "/target",
		out:             b,
		home:            path.Home(tmpHome),
		archiver:        aM,
	}
	aM.On("LoadManifestFile", path.Home(tmpHome).ArchiveDir()+"/hello-world-1.0.0.zip").Return([]byte("version: \"1.0.0\""), nil)
	aM.On("Extract", path.Home(tmpHome).ArchiveDir()+"/hello-world-1.0.0.zip", "/target", make(map[string]interface{})).Return(nil)
	err = projectCreate.run()

	aM.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "created project at \"/target\"\n", b.String())
}

func TestCreateProjectWithRegisteredTemplateAndDefinedParams(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	f := path.Home(tmpHome).TemplatesFile()
	err := ioutil.WriteFile(f, []byte(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: /my/path/new-project/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0`), 0644)
	if err != nil {
		t.Error("could not write template file")
	}

	b := bytes.NewBuffer(nil)
	params := make([]string, 2)
	params[0] = "param1=hello"
	params[1] = "param2=world"
	aM := new(ArchiverMock)
	projectCreate := &projectCreateCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		targetDir:       "/target",
		params:          params,
		out:             b,
		home:            path.Home(tmpHome),
		archiver:        aM,
	}
	aM.On("LoadManifestFile", path.Home(tmpHome).ArchiveDir()+"/hello-world-1.0.0.zip").Return([]byte(`version: "1.0.0"
parameters:
  - name: "param1"
    prompt: "Please provide a value for parameter 1"
    type: "string"
  - name: "param2"
    prompt: "Please provide a value for parameter 2"
    type: "string"`), nil)
	aM.On("Extract", path.Home(tmpHome).ArchiveDir()+"/hello-world-1.0.0.zip", "/target", map[string]interface{}{"param1": "hello", "param2": "world"}).Return(nil)
	err = projectCreate.run()

	aM.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "created project at \"/target\"\n", b.String())
}
