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

	o.MakeOptionsFromPath(config.Cfg.Directory.Model, args[0])

	o.Name = "model"

	return nil
}

// Complete completes all the required options.
func (o *ModelOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store if generating model by tables.
	if o.Table != "" {
		config.DB, _ = db.NewMySQL(config.Cfg.MysqlOptions)

		return nil
	}

	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *ModelOptions) Run(args []string) error {
	if o.Table == "" {
		return cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o.Name, o)
	}

	// Generate model from table.
	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: config.Cfg.Directory.Model,
	})

	g.UseDB(config.DB)

	g.GenerateModel(o.Table)
	g.Execute()

	return nil
}
