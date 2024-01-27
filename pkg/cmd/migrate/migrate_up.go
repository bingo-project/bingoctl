package migrate

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"
)

const (
	upUsageStr = "up"
)

var (
	upUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the up command",
		upUsageStr,
	)
)

// UpOptions is an option struct to support 'up' sub command.
type UpOptions struct {
	// Options
	*Options
}

// NewUpOptions returns an initialized UpOptions instance.
func NewUpOptions() *UpOptions {
	return &UpOptions{
		Options: opt,
	}
}

// NewCmdUp returns new initialized instance of 'up' sub command.
func NewCmdUp() *cobra.Command {
	o := NewUpOptions()

	cmd := &cobra.Command{
		Use:                   upUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Run the database migrations",
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
func (o *UpOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes all the required options.
func (o *UpOptions) Complete(cmd *cobra.Command, args []string) error {
	return err
}

// Run executes a new sub command using the specified options.
func (o *UpOptions) Run(args []string) error {
	o.Migrator().Up()

	return nil
}
