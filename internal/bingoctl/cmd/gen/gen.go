package gen

import (
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/gen"

	"github.com/bingo-project/bingoctl/config"
	"github.com/bingo-project/bingoctl/pkg/db"
	cmdutil "github.com/bingo-project/bingoctl/util"
)

const (
	genUsageStr = "gen"
)

// GenOptions is an option struct to support 'gen' sub command.
type GenOptions struct {
	TableStr string
	Tables   []string
}

// NewGenOptions returns an initialized GenOptions instance.
func NewGenOptions() *GenOptions {
	return &GenOptions{}
}

// NewCmdGen returns new initialized instance of 'gen' sub command.
func NewCmdGen() *cobra.Command {
	o := NewGenOptions()

	cmd := &cobra.Command{
		Use:                   genUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate code with db",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringVarP(&o.TableStr, "tables", "t", "", "data tables, separated by ',', example:'user,post'.")

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *GenOptions) Validate(cmd *cobra.Command, args []string) error {
	o.Tables = strings.Split(o.TableStr, ",")

	return nil
}

// Complete completes all the required options.
func (o *GenOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store
	config.DB, _ = db.NewMySQL(config.Cfg.MysqlOptions)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *GenOptions) Run(args []string) error {
	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: config.Cfg.Directory.Model,
	})

	g.UseDB(config.DB)

	for _, table := range o.Tables {
		g.GenerateModel(table)
	}

	g.Execute()

	return nil
}
