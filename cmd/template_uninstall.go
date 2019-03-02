package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

type templateUninstallCmd struct {
	out  io.Writer
	name string
	home templ.Home
}

func newTemplateUninstallCmd(out io.Writer) *cobra.Command {
	remove := &templateUninstallCmd{out: out}

	cmd := &cobra.Command{
		Use:   "uninstall [NAME]",
		Short: "uninstall a template",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("need at least one argument, name of template")
			}

			remove.home = templ.LetsGopherSettings.Home
			for i := 0; i < len(args); i++ {
				remove.name = args[i]
				if err := remove.run(); err != nil {
					return err
				}
			}
			return nil
		},
	}

	return cmd
}

func (r *templateUninstallCmd) run() error {
	templateDir := filepath.Join(r.home.TemplatesDir(), r.name)
	err := os.RemoveAll(templateDir)

	if err != nil {
		return fmt.Errorf("can't delete template directory %q", templateDir)
	}

	return removeTemplateLine(r.out, r.name, r.home)
}

func removeTemplateLine(out io.Writer, name string, home templ.Home) error {
	templatesFile := home.TemplatesFile()
	r, err := templ.LoadTemplatesFile(templatesFile)
	if err != nil {
		return err
	}

	if !r.Remove(name) {
		return fmt.Errorf("no template named %q found", name)
	}
	if err := r.WriteFile(templatesFile, 0644); err != nil {
		return err
	}

	fmt.Fprintf(out, "%q has been removed from your templates\n", name)

	return nil
}
