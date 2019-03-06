package cmd

import (
	"errors"
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type initCmd struct {
	out  io.Writer
	home templ.Home
}

func init() {
	rootCmd.AddCommand(newInitCmd(rootCmd.OutOrStderr()))
}

func newInitCmd(out io.Writer) *cobra.Command {
	i := &initCmd{out: out}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize letsgopher",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			i.home = templ.LetsGopherSettings.Home
			return i.run()
		},
	}

	return cmd
}

func (i *initCmd) run() error {
	i.createHomeDirs()
	templatesFile := i.home.TemplatesFile()
	if fi, err := os.Stat(templatesFile); err != nil {
		fmt.Fprintf(i.out, "Creating %s \n", templatesFile)
		f := templ.NewTemplatesFile()
		if err := f.WriteFile(templatesFile, 0644); err != nil {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", templatesFile)
	}
	return nil
}

func (i *initCmd) createHomeDirs() error {
	err := createDirIfNotExist(i.home.ArchiveDir())
	if err != nil {
		return err
	}

	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("could not create %s: %s", dir, err)
		}
	}
	return nil
}
