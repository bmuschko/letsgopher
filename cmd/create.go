package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/bmuschko/lets-gopher/templ/manifest"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"io"
	"os"
	"path"
	"strconv"
)

func init() {
	rootCmd.AddCommand(newCreateCmd(rootCmd.OutOrStderr()))
}

type projectCreateCmd struct {
	templateName    string
	templateVersion string
	targetDir       string
	out             io.Writer
	home            templ.Home
}

func newCreateCmd(out io.Writer) *cobra.Command {
	create := &projectCreateCmd{out: out}

	cmd := &cobra.Command{
		Use:   "create [ARGS]",
		Short: "create a new project from a template",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the template name", "the template version", "the target project directory"); err != nil {
				return err
			}

			create.templateName = args[0]
			create.templateVersion = args[1]
			create.targetDir = args[2]
			create.home = templ.LetsGopherSettings.Home
			return create.run()
		},
	}

	return cmd
}

func (a *projectCreateCmd) run() error {
	f, err := templ.LoadTemplatesFile(a.home.TemplatesFile())
	if err != nil {
		return err
	}
	if !f.Has(a.templateName, a.templateVersion) {
		return fmt.Errorf("template with name %s and version %s hasn't been installed", a.templateName, a.templateVersion)
	}

	templateName := a.templateName + "-" + a.templateVersion
	templateZIP := path.Join(a.home.ArchiveDir(), templateName+".zip")
	archiver := &templ.ZIPArchiver{}

	tb, err := archiver.LoadFile(templateZIP)
	if err != nil {
		return err
	}
	m, err := manifest.LoadManifestData(tb)
	if err != nil {
		return err
	}
	r, err := requestParameterValues(m.Parameters)
	if err != nil {
		return err
	}

	err = archiver.Extract(templateZIP, r)
	if err != nil {
		return nil
	}
	err = os.Rename(templateName, a.targetDir)
	return err
}

func requestParameterValues(params []*manifest.Parameter) (map[string]interface{}, error) {
	replacements := make(map[string]interface{})
	if len(params) > 0 {
		core.SetFancyIcons()
	}
	for _, p := range params {
		if p.Type == manifest.StringType {
			value, err := promptString(p)
			if err != nil {
				return nil, err
			}
			replacements[p.Name] = value
		} else if p.Type == manifest.IntegerType {
			value, err := promptInteger(p)
			if err != nil {
				return nil, err
			}
			replacements[p.Name] = value
		} else if p.Type == manifest.BooleanType {
			value, err := promptBoolean(p)
			if err != nil {
				return nil, err
			}
			replacements[p.Name] = value
		} else {
			return nil, fmt.Errorf("unknown parameter type %s", p.Type)
		}
	}

	return replacements, nil
}

func promptString(p *manifest.Parameter) (string, error) {
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

func promptInteger(p *manifest.Parameter) (int, error) {
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

func promptBoolean(p *manifest.Parameter) (bool, error) {
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
