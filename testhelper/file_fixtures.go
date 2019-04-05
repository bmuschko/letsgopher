package testhelper

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var tmpDirs []string

func TmpDir(t *testing.T, dir string, prefix string) string {
	t.Helper()
	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		t.Fatalf("failed to create temporary directory %s", tmpDir)
	}
	tmpDirs = append(tmpDirs, tmpDir)
	return tmpDir
}

func CleanTmpDirs(t *testing.T) {
	for _, path := range tmpDirs {
		if err := os.RemoveAll(path); err != nil {
			t.Errorf("failed to remove temporary directory %s", path)
		}
	}

	tmpDirs = make([]string, 0)
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
		os.Remove(f.Name)
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

func CreateDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}
