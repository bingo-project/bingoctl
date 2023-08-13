package make

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	controllerUsageStr = "controller NAME"
)

var (
	controllerUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the controller command",
		controllerUsageStr,
	)
)

// ControllerOptions is an option struct to support 'controller' sub command.
type ControllerOptions struct {
	*Options
}

// NewControllerOptions returns an initialized ControllerOptions instance.
func NewControllerOptions() *ControllerOptions {
	return &ControllerOptions{
		Options: opt,
	}
}

// NewCmdController returns new initialized instance of 'controller' sub command.
func NewCmdController() *cobra.Command {
	o := NewControllerOptions()

	cmd := &cobra.Command{
		Use:                   controllerUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate controller code",
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
func (o *ControllerOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, controllerUsageErrStr)
	}

	o.MakeOptionsFromPath(config.Cfg.Directory.Controller, args[0])

	o.Name = "controller"
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
func (o *ControllerOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *ControllerOptions) Run(args []string) error {
	return cmdutil.GenerateCode(o.FilePath, cmdTemplate, o.Name, o)
}
