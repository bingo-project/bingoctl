package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	bizUsageStr = "biz NAME"
)

var (
	bizUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the biz command",
		bizUsageStr,
	)
)

// BizOptions is an option struct to support 'biz' sub command.
type BizOptions struct {
	*generator.Options
}

// NewBizOptions returns an initialized BizOptions instance.
func NewBizOptions() *BizOptions {
	return &BizOptions{
		Options: opt,
	}
}

// NewCmdBiz returns new initialized instance of 'biz' sub command.
func NewCmdBiz() *cobra.Command {
	o := NewBizOptions()

	cmd := &cobra.Command{
		Use:                   bizUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate biz code",
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
func (o *BizOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, bizUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *BizOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *BizOptions) Run(args []string) error {
	return o.GenerateCode(string(generator.TmplBiz), args[0])
}
