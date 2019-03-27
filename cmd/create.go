package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/bmuschko/lets-gopher/templ/config"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(newCreateCmd(rootCmd.OutOrStderr()))
}

type projectCreateCmd struct {
	templateName    string
	templateVersion string
	targetDir       string
	params          []string
	out             io.Writer
	home            templ.Home
}

func newCreateCmd(out io.Writer) *cobra.Command {
	create := &projectCreateCmd{out: out}

	cmd := &cobra.Command{
		Use:   "create [args]",
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

	cmd.PersistentFlags().StringSliceVar(&create.params, "param", []string{}, "parameter defined as key/value pair separated by = character")
	return cmd
}

func (a *projectCreateCmd) run() error {
	f, err := config.LoadTemplatesFile(a.home.TemplatesFile())
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
	m, err := config.LoadManifestData(tb)
	if err != nil {
		return err
	}
	err = config.ValidateManifest(m)
	if err != nil {
		return err
	}
	userDefinedParams, err := mapUserDefinedParams(a.params)
	if err != nil {
		return err
	}
	r, err := requestParameterValues(userDefinedParams, m.Parameters)
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

func mapUserDefinedParams(params []string) (map[string]string, error) {
	userDefinedParams := make(map[string]string)
	for _, p := range params {
		if !strings.Contains(p, "=") {
			return nil, fmt.Errorf("user-defined parameter %s does not separate key and value by = character", p)
		}
		s := strings.Split(p, "=")
		fmt.Println(s[0])
		fmt.Println(s[1])
		userDefinedParams[s[0]] = s[1]
	}
	return userDefinedParams, nil
}

func requestParameterValues(userDefinedParams map[string]string, manifestParams []*config.Parameter) (map[string]interface{}, error) {
	replacements := make(map[string]interface{})
	if len(manifestParams) > 0 {
		core.SetFancyIcons()
	}
	for _, p := range manifestParams {
		if value, exist := userDefinedParams[p.Name]; exist {
			if p.Enum != nil && !contains(p.Enum, value) {
				return nil, fmt.Errorf("provided value '%s' is not defined in enum [%s]",
					value, strings.Join(p.Enum, ", "))
			}
			replacements[p.Name] = value
			continue
		}

		if p.Type == config.StringType {
			value, err := promptString(p)
			if err != nil {
				return nil, err
			}
			replacements[p.Name] = value
		} else if p.Type == config.IntegerType {
			value, err := promptInteger(p)
			if err != nil {
				return nil, err
			}
			replacements[p.Name] = value
		} else if p.Type == config.BooleanType {
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
