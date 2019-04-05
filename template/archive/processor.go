package archive

import "io"

// Processor replaces placeholders in text content with values.
type Processor interface {
	Process(content []byte, target io.Writer, replacements map[string]interface{}) error
}
