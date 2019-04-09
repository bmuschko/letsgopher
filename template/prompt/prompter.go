package prompt

import "github.com/bmuschko/letsgopher/template/config"

// Prompter ask for user input for a given list of parameters.
type Prompter interface {
	Prompt(p *config.Parameter, replacements map[string]interface{}) error
}
