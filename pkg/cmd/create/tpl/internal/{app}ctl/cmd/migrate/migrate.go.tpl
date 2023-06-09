package migrate

import (
	"github.com/bingo-project/component-base/cli/console"
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/apiserver/store"
)

const (
	migrateUsageStr = "migrate"
)

// MigrateOptions is an option struct to support 'migrate' sub command.
type MigrateOptions struct {
	// Options
}

// NewMigrateOptions returns an initialized MigrateOptions instance.
func NewMigrateOptions() *MigrateOptions {
	return &MigrateOptions{}
}

// NewCmdMigrate returns new initialized instance of 'migrate' sub command.
func NewCmdMigrate() *cobra.Command {
	o := NewMigrateOptions()

	cmd := &cobra.Command{
		Use:                   migrateUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Migrate database",
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
func (o *MigrateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes all the required options.
func (o *MigrateOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *MigrateOptions) Run(args []string) error {
	err := store.S.DB().AutoMigrate(
	// Migrate models here
	)
	console.ExitIf(err)

	return nil
}
