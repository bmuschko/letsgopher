package testhelper

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
)

// CreateZip creates a ZIP file for testing purposes.
func CreateZip(filename string, files []TestFile) error {
	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(file.Content))
		if err != nil {
			return err
		}
	}

	err = w.Close()
	if err != nil {
		return err
	}
	for _, f := range files {
		os.Remove(f.Name)
	}
	return nil
}

// TestFile is a text file for bundling with a ZIP file.
type TestFile struct {
	Name    string
	Content string
}

// ReadFile reads the textual content of a file.
func ReadFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
