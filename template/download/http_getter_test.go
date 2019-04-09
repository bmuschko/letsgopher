package download

import (
	"github.com/bmuschko/letsgopher/testhelper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGetForZipFile(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	zipFile := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	files := []testhelper.TestFile{
		{"file1.txt", "This is a file1"},
		{"file2.txt", "This is a file2"},
	}
	err := testhelper.CreateZip(zipFile, files)
	if err != nil {
		t.Errorf("failed to load file %s", zipFile)
	}

	h := &TestFileHandler{zipFile: zipFile}
	server := httptest.NewServer(h)
	defer server.Close()

	g := NewHTTPGetter()
	data, err := g.Get(server.URL)
	mimeType := http.DetectContentType(data.Bytes())

	assert.Nil(t, err)
	assert.Equal(t, "application/zip", mimeType)
}

type TestFileHandler struct {
	zipFile string
}

func (h *TestFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, _ := os.Open(h.zipFile)
	defer f.Close()

	b := make([]byte, 512)
	_, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(w, f)
	if err != nil {
		panic(err)
	}
}
