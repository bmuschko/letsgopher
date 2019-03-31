package download

import (
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGetForZipFile(t *testing.T) {
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
	f.Read(b)
	f.Seek(0, 0)
	io.Copy(w, f)
	return
}
