package cmd

import (
	"bytes"
	"fmt"
	"github.com/Flaque/filet"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitNonExistentHome(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := storage.Home(tmpHome).ArchiveDir()
	templatesFile := storage.Home(tmpHome).TemplatesFile()
	err := init.run()

	assert.Nil(t, err)
	assert.DirExists(t, archiveDir)
	assert.FileExists(t, templatesFile)
	assert.Equal(t, fmt.Sprintf("Creating %s \n", templatesFile), b.String())
}

func TestInitExistentHome(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := storage.Home(tmpHome).ArchiveDir()
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Errorf("failed to create directory %s", archiveDir)
	}
	templatesFile := storage.Home(tmpHome).TemplatesFile()
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
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	b := bytes.NewBuffer(nil)
	init := &initCmd{
		out:  b,
		home: storage.Home(tmpHome),
	}
	archiveDir := storage.Home(tmpHome).ArchiveDir()
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Errorf("failed to create directory %s", archiveDir)
	}
	templatesDir := storage.Home(tmpHome).TemplatesFile()
	err = os.MkdirAll(templatesDir, 0755)
	if err != nil {
		t.Errorf("failed to create directory %s", archiveDir)
	}
	err = init.run()

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%s must be a file, not a directory", templatesDir), err.Error())
	assert.Empty(t, b.String())
}
