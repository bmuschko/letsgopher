package cmd

import (
	"bytes"
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
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
		home:            templ.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
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
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	f, err := os.Create(templatesFile)
	f.WriteString(fmt.Sprintf(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: %s/archive/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0
`, tmpHome))
	defer f.Close()
	err = uninstall.run()

	result, err := ioutil.ReadFile(templatesFile)
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
		home:            templ.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
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
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	f, err := os.Create(templatesFile)
	f.WriteString(`generated: "2019-03-15T16:31:57.232715-06:00"
templates: []`)
	defer f.Close()
	err = uninstall.run()

	result, _ := ioutil.ReadFile(templatesFile)
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
		home:            templ.Home(tmpHome),
	}
	archiveDir := filepath.Join(tmpHome, "archive")
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Fatalf("failed to create archive directory %s", archiveDir)
	}
	archiveFile := fmt.Sprintf("%s/archive/hello-world-1.0.0.zip", tmpHome)
	templatesFile := filepath.Join(tmpHome, "templates.yaml")
	f, err := os.Create(templatesFile)
	f.WriteString(fmt.Sprintf(`generated: "2019-03-15T16:31:57.232715-06:00"
templates:
- archivePath: %s/archive/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0
`, archiveFile))
	defer f.Close()
	err = uninstall.run()

	result, err := ioutil.ReadFile(templatesFile)
	assert.Nil(t, err)
	testhelper.FileNotExists(t, archiveFile)
	assert.Equal(t, `generated: "2019-03-15T16:31:57.232715-06:00"
templates: []
`, string(result))
	assert.Equal(t, fmt.Sprintf("template %q has been removed\n", "hello-world"), b.String())
}
