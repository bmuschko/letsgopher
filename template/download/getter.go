package download

import "bytes"

// Getter is an interface to support GET to the specified URL.
type Getter interface {
	Get(url string) (*bytes.Buffer, error)
}
