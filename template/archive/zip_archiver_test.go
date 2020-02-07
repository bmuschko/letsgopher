package archive

import (
	"github.com/Flaque/filet"
	"github.com/bmuschko/letsgopher/testhelper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestExtractWithoutTemplateReplacement(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{Name: manifestFile, Content: "version: \"1.0.0\""},
		{Name: "file1.txt", Content: "This is a file1"},
		{Name: "file2.txt", Content: "This is a file2"},
	}
	testhelper.CreateZip(t, archive, files)
	extractedDir := filepath.Join(tmpHome, "new-project")
	manifestFile := filepath.Join(extractedDir, manifestFile)
	extractedFile1 := filepath.Join(extractedDir, "file1.txt")
	extractedFile2 := filepath.Join(extractedDir, "file2.txt")
	err := archiver.Extract(archive, extractedDir, make(map[string]interface{}))

	assert.Nil(t, err)
	assert.DirExists(t, extractedDir)
	testhelper.FileNotExists(t, manifestFile)
	assert.FileExists(t, extractedFile1)
	assert.FileExists(t, extractedFile2)

	f1 := testhelper.ReadFile(t, extractedFile1)
	assert.Equal(t, "This is a file1", f1)
	f2 := testhelper.ReadFile(t, extractedFile2)
	assert.Equal(t, "This is a file2", f2)
}

func TestExtractWithTemplateReplacement(t *testing.T) {
	t.Skip("template replacements are currently not working")

	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{Name: manifestFile, Content: "version: \"1.0.0\""},
		{Name: "file1.txt", Content: "This is a {( .a }}"},
		{Name: "file2.txt", Content: "This is a {{ .b }}"},
	}
	testhelper.CreateZip(t, archive, files)
	extractedDir := filepath.Join(tmpHome, "new-project")
	manifestFile := filepath.Join(extractedDir, manifestFile)
	extractedFile1 := filepath.Join(extractedDir, "file1.txt")
	extractedFile2 := filepath.Join(extractedDir, "file2.txt")
	replacements := make(map[string]interface{})
	replacements["a"] = "file1"
	replacements["b"] = "file2"
	err := archiver.Extract(archive, extractedDir, replacements)

	assert.Nil(t, err)
	assert.DirExists(t, extractedDir)
	testhelper.FileNotExists(t, manifestFile)
	assert.FileExists(t, extractedFile1)
	assert.FileExists(t, extractedFile2)

	f1 := testhelper.ReadFile(t, extractedFile1)
	assert.Equal(t, "This is a file1", f1)
	f2 := testhelper.ReadFile(t, extractedFile2)
	assert.Equal(t, "This is a file2", f2)
}

func TestLoadExistingManifestFile(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{Name: manifestFile, Content: "version: \"1.0.0\""},
	}
	testhelper.CreateZip(t, archive, files)
	b, err := archiver.LoadManifestFile(archive)

	assert.Nil(t, err)
	assert.Equal(t, "version: \"1.0.0\"", string(b))
}

func TestLoadNonExistentManifestFile(t *testing.T) {
	tmpHome := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{Name: "file1.txt", Content: "This is a {( .a }}"},
		{Name: "file2.txt", Content: "This is a {{ .b }}"},
	}
	testhelper.CreateZip(t, archive, files)
	b, err := archiver.LoadManifestFile(archive)

	assert.NotNil(t, err)
	assert.Equal(t, "could not locate manifest.yaml file", err.Error())
	assert.Equal(t, []byte(nil), b)
}
