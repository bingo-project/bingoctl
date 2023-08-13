package make

import (
	"fmt"

	"github.com/spf13/cobra"

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

	return nil
}

// Complete completes all the required options.
func (o *ControllerOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a new sub command using the specified options.
func (o *ControllerOptions) Run(args []string) error {
	return o.GenerateCode("controller", args[0])
}
