package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/db"
	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	modelUsageStr = "model NAME"
)

var (
	modelUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the model command",
		modelUsageStr,
	)
)

// ModelOptions is an option struct to support 'model' sub command.
type ModelOptions struct {
	*generator.Options
}

// NewModelOptions returns an initialized ModelOptions instance.
func NewModelOptions() *ModelOptions {
	return &ModelOptions{
		Options: opt,
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

	return nil
}

// Complete completes all the required options.
func (o *ModelOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store if generating model by tables.
	var err error
	if o.Table != "" {
		config.DB, err = db.NewMySQL(config.Cfg.MysqlOptions)
	}

	return err
}

// Run executes a new sub command using the specified options.
func (o *ModelOptions) Run(args []string) error {
	return o.GenerateCode(string(generator.TmplModel), args[0])
}
