package archive

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessTemplateWithoutReplacements(t *testing.T) {
	content := []byte("hello")
	buf := bytes.NewBufferString("")
	processor := TemplateProcessor{}
	replacements := make(map[string]interface{})
	replacements["var"] = "world"
	err := processor.Process(content, buf, replacements)

	assert.Nil(t, err)
	assert.Equal(t, "hello", buf.String())
}

func TestProcessTemplateWithSingleReplacement(t *testing.T) {
	content := []byte("hello {{ .var }}")
	buf := bytes.NewBufferString("")
	processor := TemplateProcessor{}
	replacements := make(map[string]interface{})
	replacements["var"] = "world"
	err := processor.Process(content, buf, replacements)

	assert.Nil(t, err)
	assert.Equal(t, "hello world", buf.String())
}

func TestProcessTemplateWithMultipleReplacements(t *testing.T) {
	content := []byte(`hello {{ .var }}
this is a test
{{ .message }}`)
	buf := bytes.NewBufferString("")
	processor := TemplateProcessor{}
	replacements := make(map[string]interface{})
	replacements["var"] = "world"
	replacements["message"] = "bye"
	err := processor.Process(content, buf, replacements)

	assert.Nil(t, err)
	assert.Equal(t, `hello world
this is a test
bye`, buf.String())
}

func TestProcessTemplateWithConditional(t *testing.T) {
	rps := []replacement{
		{"available", `
Show this section if the condition is true
`},
		{"", `
Show this section if the condition is false
`},
	}

	for _, r := range rps {
		t.Run(r.value, func(t *testing.T) {
			content := []byte(`{{ if .condition }}
Show this section if the condition is true
{{ else }}
Show this section if the condition is false
{{ end }}`)
			buf := bytes.NewBufferString("")
			processor := TemplateProcessor{}
			replacements := make(map[string]interface{})
			replacements["condition"] = r.value
			err := processor.Process(content, buf, replacements)

			assert.Nil(t, err)
			assert.Equal(t, r.expectedOutput, buf.String())
		})
	}
}

type replacement struct {
	value          string
	expectedOutput string
}
