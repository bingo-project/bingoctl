package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	crudUsageStr = "crud"
)

var (
	crudUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the crud command",
		crudUsageStr,
	)
)

// CrudOptions is an option struct to support 'crud' sub command.
type CrudOptions struct {
	*Options

	RootPackage string
	BizPath     string
	StorePath   string
	RequestPath string
	ModelPath   string
	ModelName   string
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

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *CrudOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, crudUsageErrStr)
	}

	o.RootPackage = config.Cfg.RootPackage
	o.BizPath = config.Cfg.Directory.Biz
	o.StorePath = config.Cfg.Directory.Store
	o.RequestPath = config.Cfg.Directory.Request
	o.ModelPath = config.Cfg.Directory.Model
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}

	return nil
}

// Complete completes all the required options.
func (o *CrudOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile("tpl/model.tpl")
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *CrudOptions) Run(args []string) error {
	fmt.Println("Generating files...")

	// 1.Model
	o.Name = "model"
	o.MakeOptionsFromPath(config.Cfg.Directory.Model, args[0])
	err := cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o.Name, o)
	if err != nil {
		fmt.Println("Generating model failed, err:", err)
	}

	// 2.Store
	o.Name = "store"
	o.Directory = ""
	o.PackageName = ""
	o.MakeOptionsFromPath(config.Cfg.Directory.Store, args[0])
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}
	cmdTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	err = cmdutil.GenerateGoCode(o.FilePath, string(cmdTemplateBytes), o.Name, o)
	if err != nil {
		fmt.Println("Generating store failed, err:", err)
	}

	// 3.Request
	o.Name = "request"
	o.Directory = ""
	o.PackageName = ""
	o.MakeOptionsFromPath(config.Cfg.Directory.Request, args[0])
	cmdTemplateBytes, _ = tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	err = cmdutil.GenerateGoCode(o.FilePath, string(cmdTemplateBytes), o.Name, o)
	if err != nil {
		fmt.Println("Generating request failed, err:", err)
	}

	// 4.Biz
	o.Name = "biz"
	o.Directory = ""
	o.PackageName = ""
	o.MakeOptionsFromPath(config.Cfg.Directory.Biz, args[0])
	cmdTemplateBytes, _ = tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	err = cmdutil.GenerateGoCode(o.FilePath, string(cmdTemplateBytes), o.Name, o)
	if err != nil {
		fmt.Println("Generating biz failed, err:", err)
	}

	// 5.Controller
	o.Name = "controller"
	o.Directory = ""
	o.PackageName = ""
	o.MakeOptionsFromPath(config.Cfg.Directory.Controller, args[0])
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}
	cmdTemplateBytes, _ = tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	err = cmdutil.GenerateGoCode(o.FilePath, string(cmdTemplateBytes), o.Name, o)
	if err != nil {
		fmt.Println("Generating controller failed, err:", err)
	}

	fmt.Println("done.")

	return nil
}
