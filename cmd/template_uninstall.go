package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/template/config"
	"github.com/bmuschko/lets-gopher/template/environment"
	"github.com/bmuschko/lets-gopher/template/storage"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type templateUninstallCmd struct {
	templateName    string
	templateVersion string
	out             io.Writer
	home            storage.Home
}

func newTemplateUninstallCmd(out io.Writer) *cobra.Command {
	uninstall := &templateUninstallCmd{out: out}

	cmd := &cobra.Command{
		Use:   "uninstall [name] [version]",
		Short: "uninstall a template with a given name and version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the template name", "the template version"); err != nil {
				return err
			}

			uninstall.templateName = args[0]
			uninstall.templateVersion = args[1]
			uninstall.home = environment.Settings.Home
			return uninstall.run()
		},
	}

	return cmd
}

func (r *templateUninstallCmd) run() error {
	templatesFile := r.home.TemplatesFile()
	f, err := config.LoadTemplatesFile(templatesFile)
	if err != nil {
		return err
	}

	err = deleteTemplateArchiveFile(f, r.templateName, r.templateVersion)
	if err != nil {
		return err
	}
	return removeTemplateLine(f, templatesFile, r.out, r.templateName, r.templateVersion)
}

func deleteTemplateArchiveFile(f *config.TemplatesFile, templateName string, templateVersion string) error {
	template := f.Get(templateName, templateVersion)
	if template == nil {
		return fmt.Errorf("template with name %q and version %q hasn't been installed", templateName, templateVersion)
	}

	err := os.RemoveAll(template.ArchivePath)
	if err != nil {
		return fmt.Errorf("can't delete template archive %q", template.ArchivePath)
	}

	return nil
}

func removeTemplateLine(f *config.TemplatesFile, templatesFile string, out io.Writer, templateName string, templateVersion string) error {
	if !f.Remove(templateName, templateVersion) {
		return fmt.Errorf("no template named %q found", templateName)
	}
	if err := f.WriteFile(templatesFile, 0644); err != nil {
		return err
	}

	fmt.Fprintf(out, "template %q has been removed\n", templateName)

	return nil
}
