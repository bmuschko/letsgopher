package prompt

import "github.com/bmuschko/lets-gopher/templ/config"

type Prompter interface {
	Prompt(p *config.Parameter, replacements map[string]interface{}) error
}
