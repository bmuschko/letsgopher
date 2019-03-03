package templ

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var version string

func SetVersion(v string) {
	version = v
}

type HTTPGetter struct {
	client *http.Client
}

func (g *HTTPGetter) Get(href string) (*bytes.Buffer, error) {
	return g.get(href)
}

func (g *HTTPGetter) get(href string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return buf, err
	}
	req.Header.Set("User-Agent", "letsgopher/"+strings.TrimPrefix(version, "v"))

	resp, err := g.client.Do(req)
	if err != nil {
		return buf, err
	}
	if resp.StatusCode != 200 {
		return buf, fmt.Errorf("Failed to fetch %s : %s", href, resp.Status)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = io.Copy(buf, resp.Body)
	resp.Body.Close()
	return buf, err
}

func NewHTTPGetter() *HTTPGetter {
	return &HTTPGetter{client: &http.Client{}}
}
