package testhelper

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"testing"
)

// CreateZip creates a ZIP file for testing purposes.
func CreateZip(t *testing.T, filename string, files []TestFile) {
	outFile, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create file %s. Reason: %s", filename, err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			t.Fatalf("Failed to create file %s. Reason: %s", file.Name, err)
		}
		_, err = f.Write([]byte(file.Content))
		if err != nil {
			t.Fatalf("Failed to write to file %s. Reason: %s", file.Name, err)
		}
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("Failed to close file %s. Reason: %s", outFile.Name(), err)
	}
	for _, f := range files {
		os.Remove(f.Name)
	}
}

// TestFile is a text file for bundling with a ZIP file.
type TestFile struct {
	Name    string
	Content string
}

// ReadFile reads the textual content of a file.
func ReadFile(t *testing.T, file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to read file %s. Reason: %s", file, err)
	}
	return string(b)
}

// WriteFile writes the textual content of a file.
func WriteFile(t *testing.T, file string, content string, perm os.FileMode) {
	err := ioutil.WriteFile(file, []byte(content), perm)
	if err != nil {
		t.Fatalf("Failed to write file %s. Reason: %s", file, err)
	}
}
