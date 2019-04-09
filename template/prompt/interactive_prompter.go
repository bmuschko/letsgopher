package prompt

import (
	"fmt"
	"github.com/bmuschko/letsgopher/template/config"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"strconv"
)

// InteractivePrompter ask for user input interactively on the console.
type InteractivePrompter struct {
}

func init() {
	core.SetFancyIcons()
}

// Prompt requests user input from the console.
func (ip *InteractivePrompter) Prompt(p *config.Parameter, replacements map[string]interface{}) error {
	if p.Type == config.StringType {
		value, err := promptString(p)
		if err != nil {
			return err
		}
		replacements[p.Name] = value
	} else if p.Type == config.IntegerType {
		value, err := promptInteger(p)
		if err != nil {
			return err
		}
		replacements[p.Name] = value
	} else if p.Type == config.BooleanType {
		value, err := promptBoolean(p)
		if err != nil {
			return err
		}
		replacements[p.Name] = value
	} else {
		return fmt.Errorf("unknown parameter type %s", p.Type)
	}

	return nil
}

func promptString(p *config.Parameter) (string, error) {
	value := ""
	var err error

	if p.Enum != nil {
		prompt := &survey.Select{
			Message: p.Prompt,
			Options: p.Enum,
		}
		if p.Description != "" {
			prompt.Help = p.Description
		}
		if p.DefaultValue != "" {
			prompt.Default = p.DefaultValue
		}
		err = survey.AskOne(prompt, &value, survey.Required)
	} else {
		prompt := &survey.Input{
			Message: p.Prompt,
		}
		if p.Description != "" {
			prompt.Help = p.Description
		}
		if p.DefaultValue != "" {
			prompt.Default = p.DefaultValue
		}
		err = survey.AskOne(prompt, &value, survey.Required)
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func promptInteger(p *config.Parameter) (int, error) {
	value := 0
	var err error

	if p.Enum != nil {
		prompt := &survey.Select{
			Message: p.Prompt,
			Options: p.Enum,
		}
		if p.Description != "" {
			prompt.Help = p.Description
		}
		if p.DefaultValue != "" {
			prompt.Default = p.DefaultValue
		}
		err = survey.AskOne(prompt, &value, survey.Required)
	} else {
		prompt := &survey.Input{
			Message: p.Prompt,
		}
		if p.Description != "" {
			prompt.Help = p.Description
		}
		if p.DefaultValue != "" {
			prompt.Default = p.DefaultValue
		}
		err = survey.AskOne(prompt, &value, survey.Required)
	}
	if err != nil {
		return 0, err
	}
	return value, nil
}

func promptBoolean(p *config.Parameter) (bool, error) {
	value := false
	prompt := &survey.Confirm{
		Message: p.Prompt,
	}
	if p.Description != "" {
		prompt.Help = p.Description
	}
	err := survey.AskOne(prompt, &value, survey.Required)
	if err != nil {
		return false, err
	}
	if p.DefaultValue != "" {
		b, err := strconv.ParseBool(p.DefaultValue)
		if err != nil {
			return false, err
		}
		prompt.Default = b
	}
	return value, nil
}
