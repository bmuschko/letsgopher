package cmd

import (
	"fmt"
	"github.com/bmuschko/letsgopher/template/archive"
	"github.com/bmuschko/letsgopher/template/config"
	"github.com/bmuschko/letsgopher/template/environment"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/kr/text"
	"github.com/spf13/cobra"
	"io"
	"path"
)

type templateInspectCmd struct {
	templateName    string
	templateVersion string
	out             io.Writer
	home            storage.Home
	archiver        archive.Archiver
}

func newTemplateInspectCmd(out io.Writer) *cobra.Command {
	inspect := &templateInspectCmd{out: out}

	cmd := &cobra.Command{
		Use:   "inspect [name] [version]",
		Short: "inspects a template with a given name and version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the template name", "the template version"); err != nil {
				return err
			}

			inspect.templateName = args[0]
			inspect.templateVersion = args[1]
			inspect.home = environment.Settings.Home
			inspect.archiver = &archive.ZIPArchiver{}
			return inspect.run()
		},
	}
	return cmd
}

func (a *templateInspectCmd) run() error {
	f, err := config.LoadTemplatesFile(a.home.TemplatesFile())
	if err != nil {
		return fmt.Errorf("failed to load templates.yaml file")
	}
	if !f.Has(a.templateName, a.templateVersion) {
		return fmt.Errorf("template with name %q and version %q hasn't been installed", a.templateName, a.templateVersion)
	}

	templateName := a.templateName + "-" + a.templateVersion
	templateZIP := path.Join(a.home.ArchiveDir(), templateName+".zip")

	tb, err := a.archiver.LoadManifestFile(templateZIP)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.out, "template:")
	fmt.Fprintln(a.out, fmt.Sprintf("  name: %q", a.templateName))
	fmt.Fprintln(a.out, fmt.Sprintf("  version: %q", a.templateVersion))
	fmt.Fprintln(a.out, "manifest:")
	fmt.Fprintln(a.out, text.Indent(string(tb), "  "))
	return nil
}
