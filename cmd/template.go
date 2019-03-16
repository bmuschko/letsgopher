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
		Use:   "template install|uninstall|list|inspect [args]",
		Short: "install, uninstall, list, inspect template",
	}

	cmd.AddCommand(newTemplateInstallCmd(out))
	cmd.AddCommand(newTemplateUninstallCmd(out))
	cmd.AddCommand(newTemplateListCmd(out))
	cmd.AddCommand(newTemplateInspectCmd(out))
	return cmd
}
