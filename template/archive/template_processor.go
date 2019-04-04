package archive

import (
	"io"
	"text/template"
)

type TemplateProcessor struct {
}

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
