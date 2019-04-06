package cmd

import (
	"bytes"
	"fmt"
	"github.com/bmuschko/lets-gopher/template/storage"
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUninstallExistentTemplate(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	b := bytes.NewBuffer(nil)
	uninstall := &templateUninstallCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		out:             b,
		home:            storage.Home(tmpHome),
	}
	archiveDir := storage.Home(tmpHome).ArchiveDir()
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Fatalf("failed to create archive directory %s", archiveDir)
	}
	archiveFile := filepath.Join(archiveDir, "hello-world-1.0.0.zip")
	aF, err := os.Create(archiveFile)
	if err != nil {
		t.Fatalf("failed to create archive file %s", archiveFile)
	}
	defer aF.Close()
	templatesFile := storage.Home(tmpHome).TemplatesFile()
	f, err := os.Create(templatesFile)
	if err != nil {
		t.Errorf("failed to create file %s", templatesFile)
	}
	_, err = f.WriteString(fmt.Sprintf(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: %s/archive/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0
`, tmpHome))
	if err != nil {
		t.Errorf("failed to write to file %s", f.Name())
	}
	defer f.Close()
	err = uninstall.run()

	result, e := ioutil.ReadFile(templatesFile)
	if e != nil {
		t.Errorf("failed to read file %s", templatesFile)
	}

	assert.Nil(t, err)
	testhelper.FileNotExists(t, archiveFile)
	assert.Equal(t, `generated: "2019-03-15T16:31:57.232715-06:00"
templates: []
`, string(result))
	assert.Equal(t, fmt.Sprintf("template %q has been removed\n", "hello-world"), b.String())
}

func TestUninstallNonExistentTemplate(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	b := bytes.NewBuffer(nil)
	uninstall := &templateUninstallCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		out:             b,
		home:            storage.Home(tmpHome),
	}
	archiveDir := storage.Home(tmpHome).ArchiveDir()
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Fatalf("failed to create archive directory %s", archiveDir)
	}
	archiveFile := filepath.Join(archiveDir, "web-project-1.0.0.zip")
	aF, err := os.Create(archiveFile)
	if err != nil {
		t.Fatalf("failed to create archive file %s", archiveFile)
	}
	defer aF.Close()
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
	err = uninstall.run()

	result, e := ioutil.ReadFile(templatesFile)
	if e != nil {
		t.Errorf("failed to read file %s", templatesFile)
	}

	assert.NotNil(t, err)
	assert.FileExists(t, archiveFile)
	assert.Equal(t, `generated: "2019-03-15T16:31:57.232715-06:00"
templates: []`, string(result))
	assert.Equal(t, fmt.Sprintf("template with name %q and version %q hasn't been installed", "hello-world", "1.0.0"), err.Error())
}

func TestUninstallExistentArchiveFile(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	b := bytes.NewBuffer(nil)
	uninstall := &templateUninstallCmd{
		templateName:    "hello-world",
		templateVersion: "1.0.0",
		out:             b,
		home:            storage.Home(tmpHome),
	}
	archiveDir := storage.Home(tmpHome).ArchiveDir()
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Fatalf("failed to create archive directory %s", archiveDir)
	}
	archiveFile := fmt.Sprintf("%s/archive/hello-world-1.0.0.zip", tmpHome)
	templatesFile := storage.Home(tmpHome).TemplatesFile()
	f, err := os.Create(templatesFile)
	if err != nil {
		t.Errorf("failed to create file %s", templatesFile)
	}
	_, err = f.WriteString(fmt.Sprintf(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: %s/archive/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0
`, archiveFile))
	if err != nil {
		t.Errorf("failed to write to file %s", f.Name())
	}
	defer f.Close()
	err = uninstall.run()

	result, e := ioutil.ReadFile(templatesFile)
	if e != nil {
		t.Errorf("failed to read file %s", templatesFile)
	}

	assert.Nil(t, err)
	testhelper.FileNotExists(t, archiveFile)
	assert.Equal(t, `generated: "2019-03-15T16:31:57.232715-06:00"
templates: []
`, string(result))
	assert.Equal(t, fmt.Sprintf("template %q has been removed\n", "hello-world"), b.String())
}
