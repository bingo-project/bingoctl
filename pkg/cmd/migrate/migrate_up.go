package migrate

import (
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/migrate/runner"
)

const (
	upUsageStr = "up"
)

// UpOptions is an option struct to support 'up' sub command.
type UpOptions struct {
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
			cmdutil.CheckErr(o.Run(args))
		},
	}

	return cmd
}

// Run executes a new sub command using the specified options.
func (o *UpOptions) Run(args []string) error {
	r, err := runner.NewRunner(o.Verbose, o.Rebuild)
	if err != nil {
		return err
	}

	return r.Run("up")
}
