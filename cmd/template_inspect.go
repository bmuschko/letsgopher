package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/kr/text"
	"github.com/spf13/cobra"
	"io"
	"path"
)

type templateInspectCmd struct {
	templateName    string
	templateVersion string
	out             io.Writer
	home            templ.Home
}

func newTemplateInspectCmd(out io.Writer) *cobra.Command {
	inspect := &templateInspectCmd{out: out}

	cmd := &cobra.Command{
		Use:   "inspect [NAME] [VERSION]",
		Short: "inspects a template with a given name and version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the template name", "the template version"); err != nil {
				return err
			}

			inspect.templateName = args[0]
			inspect.templateVersion = args[1]
			inspect.home = templ.LetsGopherSettings.Home
			return inspect.run()
		},
	}
	return cmd
}

func (a *templateInspectCmd) run() error {
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

	fmt.Println("template:")
	fmt.Printf("  name: \"%s\"\n", a.templateName)
	fmt.Printf("  version: \"%s\"\n", a.templateVersion)
	fmt.Println("manifest:")
	fmt.Print(text.Indent(string(tb), "  "))
	return nil
}
