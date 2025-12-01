package migrate

import (
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/migrate/runner"
)

const (
	freshUsageStr = "fresh"
)

// FreshOptions is an option struct to support 'fresh' sub command.
type FreshOptions struct {
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
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

// Run executes a new sub command using the specified options.
func (o *FreshOptions) Run(args []string) error {
	r, err := runner.NewRunner(o.Verbose, o.Rebuild)
	if err != nil {
		return err
	}

	return r.Run("fresh")
}
