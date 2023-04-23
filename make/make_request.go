package make

import (
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/config"
	cmdutil "github.com/bingo-project/bingoctl/util"
)

const (
	requestUsageStr = "request NAME"
)

// RequestOptions is an option struct to support 'request' sub command.
type RequestOptions struct {
	*Options
}

// NewRequestOptions returns an initialized RequestOptions instance.
func NewRequestOptions() *RequestOptions {
	return &RequestOptions{
		opt,
	}
}

// NewCmdRequest returns new initialized instance of 'request' sub command.
func NewCmdRequest() *cobra.Command {
	o := NewRequestOptions()

	cmd := &cobra.Command{
		Use:                   requestUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate request code",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *RequestOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, cmdUsageErrStr)
	}

	o.MakeOptionsFromPath(config.Cfg.Directory.Request, args[0])

	return nil
}

// Complete completes all the required options.
func (o *RequestOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile("tpl/request.tpl")
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *RequestOptions) Run(args []string) error {
	return cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o)
}
