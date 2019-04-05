package download

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var version string

// Sets the application version for use in HTTP request header.
func SetVersion(v string) {
	version = v
}

// HttpGetter is the default HTTP backend handler.
type HTTPGetter struct {
	client *http.Client
}

// Get performs a Get from repo.Getter and returns the body.
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
		return buf, fmt.Errorf("failed to fetch %s : %s", href, resp.Status)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = io.Copy(buf, resp.Body)
	return buf, err
}

// NewHTTPGetter constructs a valid HTTP client as HttpGetter.
func NewHTTPGetter() *HTTPGetter {
	return &HTTPGetter{client: &http.Client{}}
}
