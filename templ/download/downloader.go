package download

import (
	"github.com/bmuschko/lets-gopher/templ"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type TemplateDownloader struct {
	Home   templ.Home
	Getter Getter
}

func (c *TemplateDownloader) DownloadTo(url, dest string) (string, error) {
	data, err := c.Getter.Get(url)
	if err != nil {
		return "", err
	}

	destfile := filepath.Join(c.Home.ArchiveDir(), extractTemplateName(url))
	if err := ioutil.WriteFile(destfile, data.Bytes(), 0644); err != nil {
		return "", err
	}

	return destfile, nil
}

func extractTemplateName(url string) string {
	lastDotSlash := strings.LastIndex(url, "/")
	r := []rune(url)
	return string(r[lastDotSlash:len(url)])
}
