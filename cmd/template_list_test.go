package cmd

import (
	"bytes"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNonExistentTemplateFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	b := bytes.NewBuffer(nil)
	templateList := &templateListCmd{
		out:  b,
		home: templ.Home(tmpHome),
	}
	err = templateList.run()

	assert.NotNil(t, err)
	assert.Equal(t, "failed to load templates.yaml file", err.Error())
}

func TestEmptyTemplateList(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	b := bytes.NewBuffer(nil)
	templateList := &templateListCmd{
		out:  b,
		home: templ.Home(tmpHome),
	}
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	f, err := os.Create(templatesFile)
	f.WriteString(`generated: "2019-03-15T16:31:57.232715-06:00"
templates: []`)
	defer f.Close()
	err = templateList.run()

	assert.NotNil(t, err)
	assert.Equal(t, "no templates installed", err.Error())
}

func TestPopulatedTemplateList(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	b := bytes.NewBuffer(nil)
	templateList := &templateListCmd{
		out:  b,
		home: templ.Home(tmpHome),
	}
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	f, err := os.Create(templatesFile)
	f.WriteString(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: /Users/bmuschko/.letsgopher/archive/hello-world-0.2.0.zip
  name: hello-world
  version: 0.2.0`)
	defer f.Close()
	err = templateList.run()

	assert.Nil(t, err)
	assert.Equal(t, `NAME       	VERSION	ARCHIVE PATH                                             
hello-world	0.2.0  	/Users/bmuschko/.letsgopher/archive/hello-world-0.2.0.zip
`, b.String())
}
