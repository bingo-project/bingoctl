package make

import (
	"fmt"

	"github.com/spf13/cobra"
	"gorm.io/gen"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/db"
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
	*Options

	Table string
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

	cmd.Flags().StringVarP(&o.Table, "table", "t", "", "generate model by table, example:'post'.")

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
	if o.Table != "" {
		config.DB, _ = db.NewMySQL(config.Cfg.MysqlOptions)

		return nil
	}

	return nil
}

// Run executes a new sub command using the specified options.
func (o *ModelOptions) Run(args []string) error {
	if o.Table == "" {
		return o.GenerateCode("model", args[0])
	}

	// Generate model from table.
	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: o.Directory,

		// generate model global configuration
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(config.DB)

	// Generate struct `StructName` based on table `Table`
	g.GenerateModelAs(o.Table, o.StructName)

	g.Execute()

	return nil
}
