package make

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	middlewareUsageStr = "middleware NAME"
)

var (
	middlewareUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the middleware command",
		middlewareUsageStr,
	)
)

// MiddlewareOptions is an option struct to support 'middleware' sub command.
type MiddlewareOptions struct {
	*Options
}

// NewMiddlewareOptions returns an initialized MiddlewareOptions instance.
func NewMiddlewareOptions() *MiddlewareOptions {
	return &MiddlewareOptions{
		opt,
	}
}

// NewCmdMiddleware returns new initialized instance of 'middleware' sub command.
func NewCmdMiddleware() *cobra.Command {
	o := NewMiddlewareOptions()

	cmd := &cobra.Command{
		Use:                   middlewareUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate middleware code",
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
func (o *MiddlewareOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, middlewareUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *MiddlewareOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *MiddlewareOptions) Run(args []string) error {
	return o.GenerateCode("middleware", args[0])
}
