// ABOUTME: Database management commands for bingoctl
// ABOUTME: Parent command that groups database-related subcommands like seed
package db

import (
	"github.com/spf13/cobra"
)

var opt = NewOptions()

// Options is an option struct to support 'db' sub commands.
type Options struct {
	Verbose bool
	Rebuild bool
}

// NewOptions returns an initialized Options instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdDB returns new initialized instance of 'db' command.
func NewCmdDB() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "db COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Database management commands",
	}

	cmd.PersistentFlags().BoolVarP(&opt.Verbose, "verbose", "v", false, "Show detailed compilation output")
	cmd.PersistentFlags().BoolVar(&opt.Rebuild, "rebuild", false, "Force rebuild binary")

	cmd.AddCommand(NewCmdSeed())

	return cmd
}
