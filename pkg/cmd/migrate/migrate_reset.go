package migrate

import (
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/migrate/runner"
)

const (
	resetUsageStr = "reset"
)

// ResetOptions is an option struct to support 'reset' sub command.
type ResetOptions struct {
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
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

// Run executes a new sub command using the specified options.
func (o *ResetOptions) Run(args []string) error {
	r, err := runner.NewRunner(o.Verbose, o.Rebuild)
	if err != nil {
		return err
	}

	return r.Run("reset")
}
