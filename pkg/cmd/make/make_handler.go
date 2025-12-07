package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	handlerUsageStr = "handler NAME"
)

var (
	handlerUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the handler command",
		handlerUsageStr,
	)
)

// HandlerOptions is an option struct to support 'handler' sub command.
type HandlerOptions struct {
	*generator.Options
}

// NewHandlerOptions returns an initialized HandlerOptions instance.
func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		Options: opt,
	}
}

// NewCmdHandler returns new initialized instance of 'handler' sub command.
func NewCmdHandler() *cobra.Command {
	o := NewHandlerOptions()

	cmd := &cobra.Command{
		Use:                   handlerUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate handler code",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.ModelName, "model", "m", "", "Model name.")

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *HandlerOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, handlerUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *HandlerOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *HandlerOptions) Run(args []string) error {
	return o.GenerateCode(string(generator.TmplHandler), args[0])
}
