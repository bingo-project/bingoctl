// ABOUTME: Seed command implementation for database seeding
// ABOUTME: Uses dynamic runner to compile and execute user-defined seeders
package db

import (
	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/seed/runner"
)

// SeedOptions is an option struct to support 'seed' sub command.
type SeedOptions struct {
	*Options
	Seeder string
}

// NewSeedOptions returns an initialized SeedOptions instance.
func NewSeedOptions() *SeedOptions {
	return &SeedOptions{
		Options: opt,
	}
}

// NewCmdSeed returns new initialized instance of 'seed' sub command.
func NewCmdSeed() *cobra.Command {
	o := NewSeedOptions()

	cmd := &cobra.Command{
		Use:                   "seed",
		DisableFlagsInUseLine: true,
		Short:                 "Seed the database with records",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.Seeder, "seeder", "", "The class name of the seeder to run")

	return cmd
}

// Run executes the seed command.
func (o *SeedOptions) Run() error {
	r, err := runner.NewRunner(o.Verbose, o.Rebuild)
	if err != nil {
		return err
	}
	return r.Run(o.Seeder)
}
