package cmd

import (
	"bytes"
	"fmt"
	"github.com/bmuschko/lets-gopher/template/storage"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInitNonExistentHome(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	err = init.run()

	assert.Nil(t, err)
	assert.DirExists(t, archiveDir)
	assert.FileExists(t, templatesFile)
	assert.Equal(t, fmt.Sprintf("Creating %s \n", templatesFile), b.String())
}

func TestInitExistentHome(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	os.MkdirAll(archiveDir, 0755)
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	os.Create(templatesFile)
	err = init.run()

	assert.Nil(t, err)
	assert.DirExists(t, archiveDir)
	assert.FileExists(t, templatesFile)
	assert.Empty(t, b.String())
}

func TestInitTemplateFileIsDirectory(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	os.MkdirAll(archiveDir, 0755)
	templatesDir := filepath.Join(tmpHome, "templates.yaml")
	os.MkdirAll(templatesDir, 0755)
	err = init.run()

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s must be a file, not a directory", templatesDir), err.Error())
	assert.Empty(t, b.String())
}
