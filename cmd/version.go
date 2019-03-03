package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "templateVersion",
	Short: "print the templateVersion number and exit",
	Run:   printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func SetVersion(v string) {
	version = v
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("letsgopher %s\n", version)
}
