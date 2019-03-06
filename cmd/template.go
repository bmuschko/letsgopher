package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

func init() {
	rootCmd.AddCommand(newTemplateCmd(rootCmd.OutOrStderr()))
}

func newTemplateCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template [FLAGS] install|uninstall|list [ARGS]",
		Short: "install, uninstall, list template",
	}

	cmd.AddCommand(newTemplateInstallCmd(out))
	cmd.AddCommand(newTemplateUninstallCmd(out))
	cmd.AddCommand(newTemplateListCmd(out))
	cmd.AddCommand(newTemplateInspectCmd(out))
	return cmd
}
