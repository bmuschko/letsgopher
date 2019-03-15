package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"io"
	"os"
	"path"
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
	m, err := templ.LoadManifestData(tb)
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

func requestParameterValues(params []*templ.Parameter) (map[string]string, error) {
	replacements := make(map[string]string)
	if len(params) > 0 {
		core.SetFancyIcons()
	}
	for _, p := range params {
		value := ""
		prompt := &survey.Input{
			Message: "Please enter " + p.Description,
		}
		if p.DefaultValue != "" {
			prompt.Default = p.DefaultValue
		}
		err := survey.AskOne(prompt, &value, survey.Required)
		if err != nil {
			return nil, err
		}
		replacements[p.Name] = value
	}

	return replacements, nil
}
