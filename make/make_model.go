package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/config"
	cmdutil "github.com/bingo-project/bingoctl/util"
)

const (
	modelUsageStr = "model NAME"
)

var (
	modelUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the model command",
		bizUsageStr,
	)
)

// ModelOptions is an option struct to support 'model' sub command.
type ModelOptions struct {
	*Options
}

// NewModelOptions returns an initialized ModelOptions instance.
func NewModelOptions() *ModelOptions {
	return &ModelOptions{
		opt,
	}
}

// NewCmdModel returns new initialized instance of 'model' sub command.
func NewCmdModel() *cobra.Command {
	o := NewModelOptions()

	cmd := &cobra.Command{
		Use:                   modelUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate model code",
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
func (o *ModelOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, modelUsageErrStr)
	}

	o.MakeOptionsFromPath(config.Cfg.Directory.Model, args[0])

	return nil
}

// Complete completes all the required options.
func (o *ModelOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile("tpl/model.tpl")
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *ModelOptions) Run(args []string) error {
	return cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o)
}
