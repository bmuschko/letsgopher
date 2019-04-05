package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const manifestFile = "manifest.yaml"

// ZIPArchiver handles ZIP archive files.
type ZIPArchiver struct {
	Processor Processor
}

// Extract expands the contents of a ZIP file.
func (a *ZIPArchiver) Extract(archiveFile string, targetDir string, replacements map[string]interface{}) error {
	r, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		err := a.extractAndWriteFile(f, targetDir, replacements)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ZIPArchiver) extractAndWriteFile(f *zip.File, targetDir string, replacements map[string]interface{}) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	path := filepath.Join(targetDir, f.Name)

	if f.FileInfo().IsDir() {
		err := os.MkdirAll(path, f.Mode())
		if err != nil {
			return err
		}
	} else {
		// ignore manifest file
		if filepath.Base(path) == manifestFile {
			return nil
		}
		err := os.MkdirAll(filepath.Dir(path), f.Mode())
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		b, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil
		}
		err = a.Processor.Process(b, f, replacements)
		if err != nil {
			return nil
		}
	}
	return nil
}

// LoadManifestFile loads the manifest from a ZIP file.
func (a *ZIPArchiver) LoadManifestFile(src string) ([]byte, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		if filepath.Base(f.Name) == manifestFile {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			b := bytes.NewBuffer(nil)
			_, err = io.Copy(b, rc)
			if err != nil {
				return nil, err
			}
			return b.Bytes(), err
		}
	}
	return nil, fmt.Errorf("could not locate %s file", manifestFile)
}
