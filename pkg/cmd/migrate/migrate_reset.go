package migrate

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"
)

const (
	resetUsageStr = "reset"
)

var (
	resetUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the reset command",
		resetUsageStr,
	)
)

// ResetOptions is an option struct to support 'reset' sub command.
type ResetOptions struct {
	// Options
	*Options
}

// NewResetOptions returns an initialized ResetOptions instance.
func NewResetOptions() *ResetOptions {
	return &ResetOptions{
		Options: opt,
	}
}

// NewCmdReset returns new initialized instance of 'reset' sub command.
func NewCmdReset() *cobra.Command {
	o := NewResetOptions()

	cmd := &cobra.Command{
		Use:                   resetUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Rollback all database migrations",
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
func (o *ResetOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes all the required options.
func (o *ResetOptions) Complete(cmd *cobra.Command, args []string) error {
	return err
}

// Run executes a new sub command using the specified options.
func (o *ResetOptions) Run(args []string) error {
	o.Migrator().Reset()

	return nil
}
