package cmd

import (
	"bytes"
	"fmt"
	"github.com/bmuschko/lets-gopher/template/storage"
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestInitNonExistentHome(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	err := init.run()

	assert.Nil(t, err)
	assert.DirExists(t, archiveDir)
	assert.FileExists(t, templatesFile)
	assert.Equal(t, fmt.Sprintf("Creating %s \n", templatesFile), b.String())
}

func TestInitExistentHome(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Errorf("failed to create directory %s", archiveDir)
	}
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	_, err = os.Create(templatesFile)
	if err != nil {
		t.Errorf("failed to create file %s", templatesFile)
	}
	err = init.run()

	assert.Nil(t, err)
	assert.DirExists(t, archiveDir)
	assert.FileExists(t, templatesFile)
	assert.Empty(t, b.String())
}

func TestInitTemplateFileIsDirectory(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Errorf("failed to create directory %s", archiveDir)
	}
	templatesDir := filepath.Join(tmpHome, "templates.yaml")
	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Errorf("failed to create directory %s", archiveDir)
	}
	err = init.run()

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s must be a file, not a directory", templatesDir), err.Error())
	assert.Empty(t, b.String())
}
