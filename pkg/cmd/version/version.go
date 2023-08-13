package version

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"
)

// Options is a struct to support version command.
type Options struct{}

// NewOptions returns an initialized Options instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdVersion returns a cobra command for fetching versions.
func NewCmdVersion() *cobra.Command {
	o := NewOptions()
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the client and server version information",
		Long:  "Print the client and server version information for the current context",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

func (o *Options) Complete(cmd *cobra.Command, args []string) (err error) {
	return
}

// Validate makes sure there is no discrepancy in command options.
func (o *Options) Validate(cmd *cobra.Command, args []string) (err error) {
	return
}

// Run executes a creat subcommand using the specified options.
func (o *Options) Run(args []string) (err error) {
	fmt.Println("version: v1.0.8")

	return nil
}
