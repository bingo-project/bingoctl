package make

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"
)

const (
	jobUsageStr = "job NAME"
)

var (
	jobUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the job command",
		jobUsageStr,
	)
)

// JobOptions is an option struct to support 'job' sub command.
type JobOptions struct {
	*Options
}

// NewJobOptions returns an initialized JobOptions instance.
func NewJobOptions() *JobOptions {
	return &JobOptions{
		Options: opt,
	}
}

// NewCmdJob returns new initialized instance of 'job' sub command.
func NewCmdJob() *cobra.Command {
	o := NewJobOptions()

	cmd := &cobra.Command{
		Use:                   jobUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate job code",
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
func (o *JobOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, jobUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *JobOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *JobOptions) Run(args []string) error {
	return o.GenerateCode("job", args[0])
}
