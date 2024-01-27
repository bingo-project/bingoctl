package make

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/db"
	"github.com/bingo-project/bingoctl/pkg/generator"
)

const (
	migrationUsageStr = "migration NAME"
)

var (
	migrationUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the migration command",
		migrationUsageStr,
	)
)

// MigrationOptions is an option struct to support 'migration' sub command.
type MigrationOptions struct {
	*generator.Options
}

// NewMigrationOptions returns an initialized MigrationOptions instance.
func NewMigrationOptions() *MigrationOptions {
	return &MigrationOptions{
		Options: opt,
	}
}

// NewCmdMigration returns new initialized instance of 'migration' sub command.
func NewCmdMigration() *cobra.Command {
	o := NewMigrationOptions()

	cmd := &cobra.Command{
		Use:                   migrationUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate migration code",
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
func (o *MigrationOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, migrationUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *MigrationOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store if generating model by tables.
	var err error
	if o.Table != "" {
		config.DB, err = db.NewMySQL(config.Cfg.MysqlOptions)
	}

	return err
}

// Run executes a new sub command using the specified options.
func (o *MigrationOptions) Run(args []string) error {
	return o.GenerateCode(string(generator.TmplMigration), args[0])
}
