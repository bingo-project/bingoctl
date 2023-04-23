package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/config"
	cmdutil "github.com/bingo-project/bingoctl/util"
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
	*Options

	RootPackage string
	StorePath   string
	RequestPath string
	ModelPath   string
	ModelName   string
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
		Short:                 "A brief description of your command",
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

	o.MakeOptionsFromPath(config.Cfg.Directory.Biz, args[0])

	o.RootPackage = config.Cfg.RootPackage
	o.StorePath = config.Cfg.Directory.Store
	o.RequestPath = config.Cfg.Directory.Request
	o.ModelPath = config.Cfg.Directory.Model
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}

	return nil
}

// Complete completes all the required options.
func (o *BizOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile("tpl/biz.tpl")
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *BizOptions) Run(args []string) error {
	return cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o)
}
