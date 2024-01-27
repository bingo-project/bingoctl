package migrate

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"
)

const (
	refreshUsageStr = "refresh"
)

var (
	refreshUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the refresh command",
		refreshUsageStr,
	)
)

// RefreshOptions is an option struct to support 'refresh' sub command.
type RefreshOptions struct {
	// Options
	*Options
}

// NewRefreshOptions returns an initialized RefreshOptions instance.
func NewRefreshOptions() *RefreshOptions {
	return &RefreshOptions{
		Options: opt,
	}
}

// NewCmdRefresh returns new initialized instance of 'refresh' sub command.
func NewCmdRefresh() *cobra.Command {
	o := NewRefreshOptions()

	cmd := &cobra.Command{
		Use:                   refreshUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Reset and re-run all migrations",
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
func (o *RefreshOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes all the required options.
func (o *RefreshOptions) Complete(cmd *cobra.Command, args []string) error {
	return err
}

// Run executes a new sub command using the specified options.
func (o *RefreshOptions) Run(args []string) error {
	o.Migrator().Refresh()

	return nil
}
