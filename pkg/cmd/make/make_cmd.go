package make

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	cmdUsageStr = "cmd NAME | NAME DESCRIPTION"
)

var (
	cmdUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the cmd command",
		cmdUsageStr,
	)
)

// CmdOptions is an option struct to support 'cmd' sub command.
type CmdOptions struct {
	*Options
}

// NewCmdOptions returns an initialized CmdOptions instance.
func NewCmdOptions() *CmdOptions {
	return &CmdOptions{
		Options: opt,
	}
}

// NewCmdCMD returns new initialized instance of 'cmd' sub command.
func NewCmdCMD() *cobra.Command {
	o := NewCmdOptions()

	cmd := &cobra.Command{
		Use:                   cmdUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate cmd code",
		Long:                  "Used to generate demo command source code.",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *CmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, cmdUsageErrStr)
	}

	if len(args) > 1 {
		o.Description = args[1]
	}

	return nil
}

// Complete completes all the required options.
func (o *CmdOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *CmdOptions) Run(args []string) error {
	return o.GenerateCode("cmd", args[0])
}
