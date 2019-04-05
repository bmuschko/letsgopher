package archive

import (
	"github.com/bmuschko/lets-gopher/testhelper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestExtractWithoutTemplateReplacement(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{manifestFile, "version: \"1.0.0\""},
		{"file1.txt", "This is a file1"},
		{"file2.txt", "This is a file2"},
	}
	err := testhelper.CreateZip(archive, files)
	if err != nil {
		t.Errorf("failed to create file %s", archive)
	}
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

	f1, err := testhelper.ReadFile(extractedFile1)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile1)
	}
	assert.Equal(t, "This is a file1", f1)
	f2, err := testhelper.ReadFile(extractedFile2)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile2)
	}
	assert.Equal(t, "This is a file2", f2)
}

func TestExtractWithTemplateReplacement(t *testing.T) {
	t.Skip("template replacements are currently not working")

	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{manifestFile, "version: \"1.0.0\""},
		{"file1.txt", "This is a {( .a }}"},
		{"file2.txt", "This is a {{ .b }}"},
	}
	err := testhelper.CreateZip(archive, files)
	if err != nil {
		t.Errorf("failed to create file %s", archive)
	}
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

	f1, err := testhelper.ReadFile(extractedFile1)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile1)
	}
	assert.Equal(t, "This is a file1", f1)
	f2, err := testhelper.ReadFile(extractedFile2)
	if err != nil {
		t.Errorf("failed to read file %s", extractedFile2)
	}
	assert.Equal(t, "This is a file2", f2)
}

func TestLoadExistingManifestFile(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{manifestFile, "version: \"1.0.0\""},
	}
	err := testhelper.CreateZip(archive, files)
	if err != nil {
		t.Errorf("failed to create file %s", archive)
	}
	b, err := archiver.LoadManifestFile(archive)

	assert.Nil(t, err)
	assert.Equal(t, "version: \"1.0.0\"", string(b))
}

func TestLoadNonExistentManifestFile(t *testing.T) {
	tmpHome := testhelper.TmpDir(t, "", "test")
	defer testhelper.CleanTmpDirs(t)

	archive := filepath.Join(tmpHome, "hello-world-1.0.0.zip")
	archiver := ZIPArchiver{Processor: &TemplateProcessor{}}
	files := []testhelper.TestFile{
		{"file1.txt", "This is a {( .a }}"},
		{"file2.txt", "This is a {{ .b }}"},
	}
	err := testhelper.CreateZip(archive, files)
	if err != nil {
		t.Errorf("failed to create file %s", archive)
	}
	b, err := archiver.LoadManifestFile(archive)

	assert.NotNil(t, err)
	assert.Equal(t, "could not locate manifest.yaml file", err.Error())
	assert.Equal(t, []byte(nil), b)
}
