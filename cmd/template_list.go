package cmd

import (
	"errors"
	"fmt"
	"github.com/bmuschko/letsgopher/template/config"
	"github.com/bmuschko/letsgopher/template/environment"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"io"
)

type templateListCmd struct {
	out  io.Writer
	home storage.Home
}

func newTemplateListCmd(out io.Writer) *cobra.Command {
	list := &templateListCmd{out: out}

	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "list templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			list.home = environment.Settings.Home
			return list.run()
		},
	}

	return cmd
}

func (a *templateListCmd) run() error {
	f, err := config.LoadTemplatesFile(a.home.TemplatesFile())
	if err != nil {
		return fmt.Errorf("failed to load templates.yaml file")
	}
	if len(f.Templates) == 0 {
		return errors.New("no templates installed")
	}
	table := uitable.New()
	table.AddRow("NAME", "VERSION", "ARCHIVE PATH")
	for _, te := range f.Templates {
		table.AddRow(te.Name, te.Version, te.ArchivePath)
	}
	fmt.Fprintln(a.out, table)
	return nil
}
