package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

const globalUsage = `A simple, yet flexible project generator for Go projects.

To begin working with letsgopher, run the 'letsgopher init' command:

	$ letsgopher init

This will set up any necessary local configuration.
You will only need to run this command once.

Common actions from this point include:

- letsgopher template install:   installs a new template
- letsgopher template inspect:   inspects an already installed template
- letsgopher template list:      lists all installed templates
- letsgopher create:             creates a new project from a template

`

// NewRootCmd creates the root command of the application.
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "letsgopher",
		Short: "A simple, yet flexible project generator for Go projects.",
		Long:  globalUsage,
	}

	flags := cmd.PersistentFlags()
	_ = flags.Parse(args)
	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newInitCmd(out),
		newTemplateCmd(out),
		newCreateCmd(out),
		newVersionCmd(out),
	)

	return cmd
}

func checkArgsLength(argsReceived int, requiredArgs ...string) error {
	expectedNum := len(requiredArgs)
	if argsReceived != expectedNum {
		arg := "arguments"
		if expectedNum == 1 {
			arg = "argument"
		}
		return fmt.Errorf("this command needs %v %s: %s", expectedNum, arg, strings.Join(requiredArgs, ", "))
	}
	return nil
}
