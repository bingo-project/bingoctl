package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	storeUsageStr = "store NAME"
)

var (
	storeUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the store command",
		storeUsageStr,
	)
)

// StoreOptions is an option struct to support 'store' sub command.
type StoreOptions struct {
	*generator.Options
}

// NewStoreOptions returns an initialized StoreOptions instance.
func NewStoreOptions() *StoreOptions {
	return &StoreOptions{
		Options: opt,
	}
}

// NewCmdStore returns new initialized instance of 'store' sub command.
func NewCmdStore() *cobra.Command {
	o := NewStoreOptions()

	cmd := &cobra.Command{
		Use:                   storeUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate store code",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.ModelName, "model", "m", "", "Model name.")

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *StoreOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, storeUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *StoreOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *StoreOptions) Run(args []string) error {
	return o.GenerateCode("store", args[0])
}
