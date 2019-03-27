package download

import "bytes"

type Getter interface {
	Get(url string) (*bytes.Buffer, error)
}
