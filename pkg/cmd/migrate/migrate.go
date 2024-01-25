package migrate

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/migrate"
)

var (
	err error
)

// Options is an option struct to support 'migrate' sub command.
type Options struct {
	// Options
}

// NewOptions returns an initialized Options instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdMigrate returns new initialized instance of 'migrate' sub command.
func NewCmdMigrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "migrate COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Run the database migrations",
	}

	// Add sub commands.
	cmd.AddCommand(NewCmdUp())
	cmd.AddCommand(NewCmdRollback())
	cmd.AddCommand(NewCmdRefresh())
	cmd.AddCommand(NewCmdFresh())
	cmd.AddCommand(NewCmdReset())

	return cmd
}

func migrator() *migrate.Migrator {
	return migrate.NewMigrator(strings.TrimRight(config.Cfg.Directory.Migration, "/") + "/")
}
