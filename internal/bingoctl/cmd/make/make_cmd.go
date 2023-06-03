package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/config"
	cmdutil "github.com/bingo-project/bingoctl/util"
)

const (
	cmdUsageStr = "cmd NAME | NAME DESCRIPTION"
)

var (
	cmdUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the cmd command",
		cmdUsageStr,
	)
	// tpl.
	cmdTemplate string
)

// CmdOptions is an option struct to support 'cmd' sub command.
type CmdOptions struct {
	*Options

	// Command template options
	CommandDescription string
}

// NewCmdOptions returns an initialized CmdOptions instance.
func NewCmdOptions() *CmdOptions {
	return &CmdOptions{
		Options:            opt,
		CommandDescription: "A brief description of your command",
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
		o.CommandDescription = args[1]
	}

	o.MakeOptionsFromPath(config.Cfg.Directory.CMD, args[0])

	o.Name = "cmd"

	return nil
}

// Complete completes all the required options.
func (o *CmdOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *CmdOptions) Run(args []string) error {
	return cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o.Name, o)
}
