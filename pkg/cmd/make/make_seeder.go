package make

import (
	"fmt"

	cmdutil "github.com/bingo-project/component-base/cli/util"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/generator"
)

const (
	seederUsageStr = "seeder"
)

var (
	seederUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the seeder command",
		seederUsageStr,
	)
)

// SeedOptions is an option struct to support 'seeder' sub command.
type SeedOptions struct {
	*generator.Options
}

// NewSeederOptions returns an initialized SeedOptions instance.
func NewSeederOptions() *SeedOptions {
	return &SeedOptions{
		Options: opt,
	}
}

// NewCmdSeeder returns new initialized instance of 'seeder' sub command.
func NewCmdSeeder() *cobra.Command {
	o := NewSeederOptions()

	cmd := &cobra.Command{
		Use:                   seederUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate seeder code",
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
func (o *SeedOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, seederUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *SeedOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *SeedOptions) Run(args []string) error {
	return o.GenerateCode(string(generator.TmplSeeder), args[0])
}
