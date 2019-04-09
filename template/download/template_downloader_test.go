package download

import (
	"bytes"
	"errors"
	"github.com/Flaque/filet"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/bmuschko/letsgopher/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadSuccessfully(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	zipFile := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	files := []testhelper.TestFile{
		{"file1.txt", "This is a file1"},
		{"file2.txt", "This is a file2"},
	}
	err := testhelper.CreateZip(zipFile, files)
	if err != nil {
		t.Errorf("failed to create file %s", zipFile)
	}

	gM := new(GetterMock)
	downloader := &TemplateDownloader{Getter: gM, Home: storage.Home(tmpHome)}
	url := "https://dl.dropboxusercontent.com/s/002j89do6epotqs/hello-world-1.0.0.zip"
	targetDir := storage.Home(tmpHome).ArchiveDir()
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		t.Errorf("failed to create directory %s", targetDir)
	}
	buf := bytes.NewBuffer(nil)
	source, err := os.Open(zipFile)
	if err != nil {
		t.Errorf("failed to open file %s", zipFile)
	}
	_, err = io.Copy(buf, source)
	if err != nil {
		t.Errorf("failed to copy file %s", source.Name())
	}
	destfile := filepath.Join(targetDir, "hello-world-1.0.0.zip")
	gM.On("Get", url).Return(buf, nil)
	d, err := downloader.Download(url)

	assert.Nil(t, err)
	assert.Equal(t, destfile, d)
}

func TestDownloadFailed(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	gM := new(GetterMock)
	downloader := &TemplateDownloader{Getter: gM, Home: storage.Home(tmpHome)}
	url := "https://dl.dropboxusercontent.com/s/002j89do6epotqs/hello-world-1.0.0.zip"
	downloadError := errors.New("expected")
	gM.On("Get", url).Return(bytes.NewBuffer(nil), downloadError)
	d, err := downloader.Download(url)

	assert.NotNil(t, err)
	assert.Equal(t, "expected", err.Error())
	assert.Equal(t, "", d)
}

type GetterMock struct {
	mock.Mock
}

func (g *GetterMock) Get(url string) (*bytes.Buffer, error) {
	args := g.Called(url)
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}
