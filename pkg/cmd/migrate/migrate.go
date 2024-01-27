package migrate

import (
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/bingo-project/bingoctl/pkg/migrate"
)

var (
	opt = NewOptions()
	err error
)

// Options is an option struct to support 'migrate' sub command.
type Options struct {
	DB *gorm.DB
}

// NewOptions returns an initialized Options instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdMigrate returns new initialized instance of 'migrate' sub command.
func NewCmdMigrate(db *gorm.DB) *cobra.Command {
	opt.DB = db

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

func (o *Options) Migrator() *migrate.Migrator {
	return migrate.NewMigrator(o.DB)
}
