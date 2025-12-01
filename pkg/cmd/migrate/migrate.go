package migrate

import (
	"errors"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/bingo-project/bingoctl/pkg/migrate"
)

var (
	opt = NewOptions()
	err error

	ErrInProduction = errors.New("application in production, use --force or -f to confirm")
)

// Options is an option struct to support 'migrate' sub command.
type Options struct {
	DB         *gorm.DB
	Production bool
	Force      bool
	Verbose    bool
	Rebuild    bool
}

// NewOptions returns an initialized Options instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdMigrate returns new initialized instance of 'migrate' sub command.
func NewCmdMigrate(db *gorm.DB, production bool) *cobra.Command {
	opt.DB = db
	opt.Production = production

	cmd := &cobra.Command{
		Use:                   "migrate COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Run the database migrations",
	}

	cmd.PersistentFlags().BoolVarP(&opt.Force, "force", "f", false, "Force run migration command in production")
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

func (o *Options) Migrator() *migrate.Migrator {
	return migrate.NewMigrator(o.DB)
}
