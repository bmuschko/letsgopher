package archive

import (
	"io"
	"text/template"
)

// TemplateProcessor replaces placeholders in text content with values using Go's templating functionality.
type TemplateProcessor struct {
}

// Process performs placeholder replacement.
func (tp *TemplateProcessor) Process(content []byte, target io.Writer, replacements map[string]interface{}) error {
	template, err := template.New("").Parse(string(content))
	if err != nil {
		return nil
	}
	err = template.ExecuteTemplate(target, "", replacements)
	if err != nil {
		return nil
	}
	return nil
}
