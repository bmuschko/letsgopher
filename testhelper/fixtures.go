package testhelper

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TmpDir(t *testing.T, dir string, prefix string) string {
	t.Helper()
	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpDir)
	}
	defer os.RemoveAll(tmpDir)
	return tmpDir
}

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
		err = os.Remove(f.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

type TestFile struct {
	Name    string
	Content string
}

func ReadFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
