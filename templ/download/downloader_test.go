package download

import (
	"bytes"
	"errors"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadSuccessfully(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	zipFile := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	files := []testhelper.TestFile{
		{"file1.txt", "This is a file1"},
		{"file2.txt", "This is a file2"},
	}
	testhelper.CreateZip(zipFile, files)

	gM := new(GetterMock)
	downloader := &TemplateDownloader{Getter: gM, Home: templ.Home(tmpHome)}
	url := "https://dl.dropboxusercontent.com/s/002j89do6epotqs/hello-world-1.0.0.zip"
	targetDir := filepath.Join(tmpHome, "archive")
	err = os.MkdirAll(targetDir, os.ModePerm)
	buf := bytes.NewBuffer(nil)
	source, err := os.Open(zipFile)
	_, err = io.Copy(buf, source)
	destfile := filepath.Join(targetDir, "hello-world-1.0.0.zip")
	gM.On("Get", url).Return(buf, nil)
	d, err := downloader.Download(url)

	assert.Nil(t, err)
	assert.Equal(t, destfile, d)
}

func TestDownloadFailed(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	gM := new(GetterMock)
	downloader := &TemplateDownloader{Getter: gM, Home: templ.Home(tmpHome)}
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
