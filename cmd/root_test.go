package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"testing"
)

func TestRootCmd(t *testing.T) {
	args := make([]string, 0)
	cmd := NewRootCmd(args)
	cmd.SetOutput(ioutil.Discard)
	cmd.SetArgs(args)
	cmd.Run = func(*cobra.Command, []string) {}
	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
