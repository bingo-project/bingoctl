package migrate

import (
	"fmt"

	"github.com/bingo-project/component-base/cli/console"
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"
)

const (
	freshUsageStr = "fresh"
)

var (
	freshUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the fresh command",
		freshUsageStr,
	)
)

// FreshOptions is an option struct to support 'fresh' sub command.
type FreshOptions struct {
	// Options
	*Options
}

// NewFreshOptions returns an initialized FreshOptions instance.
func NewFreshOptions() *FreshOptions {
	return &FreshOptions{
		Options: opt,
	}
}

// NewCmdFresh returns new initialized instance of 'fresh' sub command.
func NewCmdFresh() *cobra.Command {
	o := NewFreshOptions()

	cmd := &cobra.Command{
		Use:                   freshUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Drop all tables and re-run all migrations",
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
func (o *FreshOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.Production && !o.Force {
		console.Exit(ErrInProduction.Error())
	}

	return nil
}

// Complete completes all the required options.
func (o *FreshOptions) Complete(cmd *cobra.Command, args []string) error {
	return err
}

// Run executes a new sub command using the specified options.
func (o *FreshOptions) Run(args []string) error {
	o.Migrator().Fresh()

	return nil
}
