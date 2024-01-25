package migrate

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/db"
)

const (
	rollbackUsageStr = "rollback"
)

var (
	rollbackUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the rollback command",
		rollbackUsageStr,
	)
)

// RollbackOptions is an option struct to support 'rollback' sub command.
type RollbackOptions struct {
	// Options
}

// NewRollbackOptions returns an initialized RollbackOptions instance.
func NewRollbackOptions() *RollbackOptions {
	return &RollbackOptions{}
}

// NewCmdRollback returns new initialized instance of 'rollback' sub command.
func NewCmdRollback() *cobra.Command {
	o := NewRollbackOptions()

	cmd := &cobra.Command{
		Use:                   rollbackUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Rollback the last database migration",
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
func (o *RollbackOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes all the required options.
func (o *RollbackOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store
	config.DB, err = db.NewMySQL(config.Cfg.MysqlOptions)

	return err
}

// Run executes a new sub command using the specified options.
func (o *RollbackOptions) Run(args []string) error {
	migrator().Rollback()

	return nil
}
