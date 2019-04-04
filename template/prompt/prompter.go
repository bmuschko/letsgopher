package prompt

import "github.com/bmuschko/lets-gopher/template/config"

type Prompter interface {
	Prompt(p *config.Parameter, replacements map[string]interface{}) error
}
