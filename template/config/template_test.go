package config

import (
	"github.com/Flaque/filet"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/bmuschko/letsgopher/testhelper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestNewTemplatesFile(t *testing.T) {
	templatesFile := NewTemplatesFile()

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.Exactly(t, templatesFile.Templates, []*Template{})
}

func TestAddTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	webProject := &Template{Name: "web-project", Version: "2.4.1", ArchivePath: "/my/path/archive/web-project-2.4.1.zip"}
	templatesFile.Add(helloWorldTemplate)
	templatesFile.Add(webProject)

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate, webProject})
}

func TestRemoveRegisteredTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	templatesFile.Add(helloWorldTemplate)
	result := templatesFile.Remove("hello-world", "1.0.0")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.True(t, result)
	assert.Exactly(t, templatesFile.Templates, []*Template{})
}

func TestRemoveUnknownTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	result := templatesFile.Remove("hello-world", "1.0.0")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.False(t, result)
	assert.Exactly(t, templatesFile.Templates, []*Template{})
}

func TestContainsRegisteredTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	templatesFile.Add(helloWorldTemplate)
	result := templatesFile.Has("hello-world", "1.0.0")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.True(t, result)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate})
}

func TestDoesNotContainUnknownTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	templatesFile.Add(helloWorldTemplate)
	result := templatesFile.Has("web-project", "2.4.1")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.False(t, result)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate})
}

func TestGetRegisteredTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	templatesFile.Add(helloWorldTemplate)
	found := templatesFile.Get("hello-world", "1.0.0")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.Equal(t, helloWorldTemplate, found)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate})
}

func TestGetUnknownTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	templatesFile.Add(helloWorldTemplate)
	found := templatesFile.Get("web-project", "2.4.1")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.Nil(t, found)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate})
}

func TestUpdateRegisteredTemplate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	templatesFile.Add(helloWorldTemplate)
	helloWorldTemplate.ArchivePath = "/other/path/archive/hello-world-1.0.0.zip"
	result := templatesFile.Update(helloWorldTemplate)
	updatedTemplate := templatesFile.Get("hello-world", "1.0.0")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.True(t, result)
	assert.Equal(t, helloWorldTemplate, updatedTemplate)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate})
}

func TestInsertsUnknownTemplateOnUpdate(t *testing.T) {
	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	helloWorldTemplate.ArchivePath = "/other/path/archive/hello-world-1.0.0.zip"
	result := templatesFile.Update(helloWorldTemplate)
	updatedTemplate := templatesFile.Get("hello-world", "1.0.0")

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.False(t, result)
	assert.Equal(t, helloWorldTemplate, updatedTemplate)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate})
}

func TestWriteTemplateFile(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	templatesFile := NewTemplatesFile()
	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	webProject := &Template{Name: "web-project", Version: "2.4.1", ArchivePath: "/my/path/archive/web-project-2.4.1.zip"}
	templatesFile.Add(helloWorldTemplate)
	templatesFile.Add(webProject)
	f := filepath.Join(tmpHome, "template.yaml")
	err := templatesFile.WriteFile(f, 0644)
	if err != nil {
		t.Error("could not write template file")
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		t.Error("could not read template file")
	}

	assert.Nil(t, err)
	assert.Regexp(t, `generated: ".*"
templates:
- archivePath: /my/path/archive/hello-world-1.0.0.zip
  name: hello-world
  version: 1.0.0
- archivePath: /my/path/archive/web-project-2.4.1.zip
  name: web-project
  version: 2.4.1
`, string(b))
}

func TestReadTemplateFile(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	f := storage.Home(tmpHome).TemplatesFile()
	testhelper.WriteFile(t, f, `generated: "2019-03-21T08:49:27.10175-06:00"
templates:
  - archivePath: /my/path/archive/hello-world-1.0.0.zip
    name: hello-world
    version: 1.0.0
  - archivePath: /my/path/archive/web-project-2.4.1.zip
    name: web-project
    version: 2.4.1`, 0644)
	templatesFile, err := LoadTemplatesFile(f)
	if err != nil {
		t.Error("could not load template file")
	}

	helloWorldTemplate := &Template{Name: "hello-world", Version: "1.0.0", ArchivePath: "/my/path/archive/hello-world-1.0.0.zip"}
	webProject := &Template{Name: "web-project", Version: "2.4.1", ArchivePath: "/my/path/archive/web-project-2.4.1.zip"}

	assert.NotNil(t, templatesFile)
	assert.NotNil(t, templatesFile.Generated)
	assert.Exactly(t, templatesFile.Templates, []*Template{helloWorldTemplate, webProject})
}
