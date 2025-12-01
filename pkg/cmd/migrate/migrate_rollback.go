package migrate

import (
	"github.com/bingo-project/component-base/cli/console"
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/migrate/runner"
)

const (
	rollbackUsageStr = "rollback"
)

// RollbackOptions is an option struct to support 'rollback' sub command.
type RollbackOptions struct {
	// Options
	*Options
}

// NewRollbackOptions returns an initialized RollbackOptions instance.
func NewRollbackOptions() *RollbackOptions {
	return &RollbackOptions{
		Options: opt,
	}
}

// NewCmdRollback returns new initialized instance of 'rollback' sub command.
func NewCmdRollback() *cobra.Command {
	o := NewRollbackOptions()

	cmd := &cobra.Command{
		Use:                   rollbackUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Rollback the last database migration",
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
func (o *RollbackOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.Production && !o.Force {
		console.Exit(ErrInProduction.Error())
	}

	return nil
}

// Complete completes all the required options.
func (o *RollbackOptions) Complete(cmd *cobra.Command, args []string) error {
	return err
}

// Run executes a new sub command using the specified options.
func (o *RollbackOptions) Run(args []string) error {
	r, err := runner.NewRunner(o.Verbose, o.Rebuild)
	if err != nil {
		return err
	}

	return r.Run("rollback")
}
