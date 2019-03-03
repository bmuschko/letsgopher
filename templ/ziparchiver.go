package templ

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type ZIPArchiver struct {
}

func (a *ZIPArchiver) Extract(src string, replacements map[string]string) error {
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
			os.MkdirAll(f.Name, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(f.Name), f.Mode())
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
			tmpl, err := template.New(f.Name()).Parse(string(b))
			if err != nil {
				return nil
			}
			err = tmpl.ExecuteTemplate(f, f.Name(), replacements)
			if err != nil {
				return nil
			}
			//_, err = io.Copy(f, rc)
			//if err != nil {
			//	return err
			//}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	//templateDir := filepath.Join(a.Home.TemplatesDir(), dest)
	//err := archiver.Unarchive(src, templateDir)
	//if err != nil {
	//	return "", nil
	//}
	return nil
}

func (a *ZIPArchiver) LoadFile(src string) ([]byte, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		if filepath.Base(f.Name) == "manifest.yaml" {
			rc, err := f.Open()
			if err != nil {
				return nil, nil
			}
			b := bytes.NewBuffer(nil)
			_, err = io.Copy(b, rc)
			return b.Bytes(), err
		}
	}
	return nil, errors.New("could not locate manifest.yaml file")
}
