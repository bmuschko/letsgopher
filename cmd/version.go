package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
)

const appLabel = "letsgopher"

var version string

func SetVersion(v string) {
	version = v
}

type versionCmd struct {
	out io.Writer
}

func newVersionCmd(out io.Writer) *cobra.Command {
	version := &versionCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version number and exit",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return version.run()
		},
	}
	return cmd
}

func (v *versionCmd) run() error {
	_, err := fmt.Fprintf(v.out, "%s %s\n", appLabel, version)
	if err != nil {
		return err
	}
	return nil
}
