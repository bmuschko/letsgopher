package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	manifestFile = "manifest.yaml"
)

type ZIPArchiver struct {
}

func (a *ZIPArchiver) Extract(src string, replacements map[string]interface{}) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(f.Name, f.Mode())
			if err != nil {
				return err
			}
		} else {
			if filepath.Base(f.Name) == manifestFile {
				return nil
			}
			err := os.MkdirAll(filepath.Dir(f.Name), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
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
			err = processAsTemplate(b, f, replacements)
			if err != nil {
				return nil
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func processAsTemplate(b []byte, f *os.File, replacements map[string]interface{}) error {
	tmpl, err := template.New(f.Name()).Parse(string(b))
	if err != nil {
		return nil
	}
	err = tmpl.ExecuteTemplate(f, f.Name(), replacements)
	if err != nil {
		return nil
	}
	return nil
}

func (a *ZIPArchiver) LoadFile(src string) ([]byte, error) {
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
