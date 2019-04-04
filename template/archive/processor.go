package archive

import "io"

type Processor interface {
	Process(content []byte, target io.Writer, replacements map[string]interface{}) error
}
