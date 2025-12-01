package migrate

import (
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/migrate/runner"
)

const (
	refreshUsageStr = "refresh"
)

// RefreshOptions is an option struct to support 'refresh' sub command.
type RefreshOptions struct {
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
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

// Run executes a new sub command using the specified options.
func (o *RefreshOptions) Run(args []string) error {
	r, err := runner.NewRunner(o.Verbose, o.Rebuild)
	if err != nil {
		return err
	}

	return r.Run("refresh")
}
