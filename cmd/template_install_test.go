package cmd

import (
	"bytes"
	"errors"
	"github.com/Flaque/filet"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/bmuschko/letsgopher/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInstallNewTemplate(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	f := storage.Home(tmpHome).TemplatesFile()
	testhelper.WriteFile(t, f, `generated: "2019-03-21T08:49:27.10175-06:00"
templates: []`, 0644)

	b := bytes.NewBuffer(nil)
	dM := new(DownloaderMock)
	templateInstall := &templateInstallCmd{
		templateURL:  "http://my.repo.com/hello-world-1.0.0.zip",
		templateName: "new-project",
		out:          b,
		home:         storage.Home(tmpHome),
		downloader:   dM,
	}
	dM.On("Download", "http://my.repo.com/hello-world-1.0.0.zip").Return("/my/path/new-project/hello-world-1.0.0.zip", nil)
	err := templateInstall.run()

	templates := testhelper.ReadFile(t, f)

	dM.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "\"new-project\" has been added to your templates\n", b.String())
	assert.Equal(t, `generated: "2019-03-21T08:49:27.10175-06:00"
templates:
- archivePath: /my/path/new-project/hello-world-1.0.0.zip
  name: new-project
  version: 1.0.0
`, templates)
}

func TestInstallExistingTemplate(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	f := storage.Home(tmpHome).TemplatesFile()
	testhelper.WriteFile(t, f, `generated: "2019-03-21T08:49:27.10175-06:00"
templates:
- archivePath: /my/path/new-project/hello-world-1.0.0.zip
  name: new-project
  version: 1.0.0
`, 0644)

	b := bytes.NewBuffer(nil)
	dM := new(DownloaderMock)
	templateInstall := &templateInstallCmd{
		templateURL:  "http://my.repo.com/hello-world-1.0.0.zip",
		templateName: "new-project",
		out:          b,
		home:         storage.Home(tmpHome),
		downloader:   dM,
	}
	dM.On("Download", "http://my.repo.com/hello-world-1.0.0.zip").Return("/my/path/new-project/hello-world-1.0.0.zip", nil)
	err := templateInstall.run()

	dM.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, "template with name \"new-project\" already exists, please specify a different name", err.Error())
}

func TestInstallForFailedTemplateDownload(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	templatesContent := `generated: "2019-03-21T08:49:27.10175-06:00"
templates: []`
	f := storage.Home(tmpHome).TemplatesFile()
	testhelper.WriteFile(t, f, templatesContent, 0644)

	b := bytes.NewBuffer(nil)
	dM := new(DownloaderMock)
	templateInstall := &templateInstallCmd{
		templateURL:  "http://my.repo.com/hello-world-1.0.0.zip",
		templateName: "new-project",
		out:          b,
		home:         storage.Home(tmpHome),
		downloader:   dM,
	}
	dM.On("Download", "http://my.repo.com/hello-world-1.0.0.zip").Return("", errors.New("expected"))
	err := templateInstall.run()

	dM.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, "expected", err.Error())
}

type DownloaderMock struct {
	mock.Mock
}

func (d *DownloaderMock) Download(url string) (string, error) {
	args := d.Called(url)
	return args.String(0), args.Error(1)
}
