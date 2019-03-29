package archive

import (
	"archive/zip"
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractWithoutTemplateReplacement(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testFile{
		{manifestFile, "version: \"1.0.0\""},
		{"file1.txt", "This is a file1"},
		{"file2.txt", "This is a file2"},
	}
	createZip(archive, files)
	extractedDir := filepath.Join(tmpHome, "new-project")
	manifestFile := filepath.Join(extractedDir, manifestFile)
	extractedFile1 := filepath.Join(extractedDir, "file1.txt")
	extractedFile2 := filepath.Join(extractedDir, "file2.txt")
	err = archiver.Extract(archive, extractedDir, make(map[string]interface{}))

	assert.Nil(t, err)
	assert.DirExists(t, extractedDir)
	testhelper.FileNotExists(t, manifestFile)
	assert.FileExists(t, extractedFile1)
	assert.FileExists(t, extractedFile2)

	f1, err := readFile(extractedFile1)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile1)
	}
	assert.Equal(t, "This is a file1", f1)
	f2, err := readFile(extractedFile2)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile2)
	}
	assert.Equal(t, "This is a file2", f2)
}

func TestExtractWithTemplateReplacement(t *testing.T) {
	t.Skip("template replacements are currently not working")

	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testFile{
		{manifestFile, "version: \"1.0.0\""},
		{"file1.txt", "This is a {( .a }}"},
		{"file2.txt", "This is a {{ .b }}"},
	}
	createZip(archive, files)
	extractedDir := filepath.Join(tmpHome, "new-project")
	manifestFile := filepath.Join(extractedDir, manifestFile)
	extractedFile1 := filepath.Join(extractedDir, "file1.txt")
	extractedFile2 := filepath.Join(extractedDir, "file2.txt")
	replacements := make(map[string]interface{})
	replacements["a"] = "file1"
	replacements["b"] = "file2"
	err = archiver.Extract(archive, extractedDir, replacements)

	assert.Nil(t, err)
	assert.DirExists(t, extractedDir)
	testhelper.FileNotExists(t, manifestFile)
	assert.FileExists(t, extractedFile1)
	assert.FileExists(t, extractedFile2)

	f1, err := readFile(extractedFile1)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile1)
	}
	assert.Equal(t, "This is a file1", f1)
	f2, err := readFile(extractedFile2)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile2)
	}
	assert.Equal(t, "This is a file2", f2)
}

func TestLoadExistingManifestFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testFile{
		{manifestFile, "version: \"1.0.0\""},
	}
	createZip(archive, files)
	b, err := archiver.LoadManifestFile(archive)

	assert.Nil(t, err)
	assert.Equal(t, "version: \"1.0.0\"", string(b))
}

func TestLoadNonExistentManifestFile(t *testing.T) {
	tmpHome, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpHome)
	}
	defer os.RemoveAll(tmpHome)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testFile{
		{"file1.txt", "This is a {( .a }}"},
		{"file2.txt", "This is a {{ .b }}"},
	}
	createZip(archive, files)
	b, err := archiver.LoadManifestFile(archive)

	assert.NotNil(t, err)
	assert.Equal(t, "could not locate manifest.yaml file", err.Error())
	assert.Equal(t, []byte(nil), b)
}

func createZip(filename string, files []testFile) error {
	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	for _, file := range files {
		f, err := w.Create(file.name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(file.content))
		if err != nil {
			return err
		}
	}

	err = w.Close()
	if err != nil {
		return err
	}
	for _, f := range files {
		err = os.Remove(f.name)
		if err != nil {
			return err
		}
	}
	return nil
}

type testFile struct {
	name    string
	content string
}

func readFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
