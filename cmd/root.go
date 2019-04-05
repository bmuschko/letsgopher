package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// NewRootCmd creates the root command of the application.
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "letsgopher",
		Short: "Letsgopher is a project generator for Golang projects",
		Long:  "A flexible and customizable project generator for Golang projects.",
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
