package make

import (
	"fmt"

	"github.com/bingo-project/component-base/cli/console"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/db"
	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	crudUsageStr = "crud NAME"
)

var (
	crudUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the crud command",
		crudUsageStr,
	)
)

// CrudOptions is an option struct to support 'crud' sub command.
type CrudOptions struct {
	*generator.Options
}

// NewCrudOptions returns an initialized CrudOptions instance.
func NewCrudOptions() *CrudOptions {
	return &CrudOptions{
		Options: opt,
	}
}

// NewCmdCrud returns new initialized instance of 'crud' sub command.
func NewCmdCrud() *cobra.Command {
	o := NewCrudOptions()

	cmd := &cobra.Command{
		Use:                   crudUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate crud code",
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
func (o *CrudOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, crudUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *CrudOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store if generating model by tables.
	var err error
	if o.Table != "" {
		config.DB, err = db.NewMySQL(config.Cfg.MysqlOptions)
	}

	return err
}

// Run executes a new sub command using the specified options.
func (o *CrudOptions) Run(args []string) error {
	// 1.Model
	err := o.GenerateCode(string(generator.TmplModel), args[0])
	if err != nil {
		console.Error(err.Error())
	}

	// 2.Store
	o.ReSetDirectory()
	err = o.GenerateCode(string(generator.TmplStore), args[0])
	if err != nil {
		console.Error(err.Error())
	}

	// 3.Request
	o.ReSetDirectory()
	err = o.GenerateCode(string(generator.TmplRequest), args[0])
	if err != nil {
		console.Error(err.Error())
	}

	// 4.Biz
	o.ReSetDirectory()
	err = o.GenerateCode(string(generator.TmplBiz), args[0])
	if err != nil {
		console.Error(err.Error())
	}

	// 5.Controller
	o.ReSetDirectory()
	err = o.GenerateCode(string(generator.TmplController), args[0])
	if err != nil {
		console.Error(err.Error())
	}

	fmt.Println("done.")

	return nil
}
