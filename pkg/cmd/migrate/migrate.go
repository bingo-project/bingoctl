package migrate

import (
	"github.com/spf13/cobra"
)

var opt = NewOptions()

// Options is an option struct to support 'migrate' sub command.
type Options struct {
	Verbose bool
	Rebuild bool
}

// NewOptions returns an initialized Options instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdMigrate returns migrate command that uses dynamic runner.
func NewCmdMigrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "migrate COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Run the database migrations",
	}

	cmd.PersistentFlags().BoolVarP(&opt.Verbose, "verbose", "v", false, "Show detailed compilation output")
	cmd.PersistentFlags().BoolVar(&opt.Rebuild, "rebuild", false, "Force rebuild migration binary")

	// Add sub commands.
	cmd.AddCommand(NewCmdUp())
	cmd.AddCommand(NewCmdRollback())
	cmd.AddCommand(NewCmdRefresh())
	cmd.AddCommand(NewCmdFresh())
	cmd.AddCommand(NewCmdReset())

	return cmd
}
