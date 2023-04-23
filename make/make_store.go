package make

import (
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/config"
	cmdutil "github.com/bingo-project/bingoctl/util"
)

const (
	storeUsageStr = "store NAME"
)

// StoreOptions is an option struct to support 'store' sub command.
type StoreOptions struct {
	*Options

	RootPackage string
	ModelPath   string
	ModelName   string
}

// NewStoreOptions returns an initialized StoreOptions instance.
func NewStoreOptions() *StoreOptions {
	return &StoreOptions{
		Options: opt,
	}
}

// NewCmdStore returns new initialized instance of 'store' sub command.
func NewCmdStore() *cobra.Command {
	o := NewStoreOptions()

	cmd := &cobra.Command{
		Use:                   storeUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate store code",
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
func (o *StoreOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, cmdUsageErrStr)
	}

	o.MakeOptionsFromPath(config.Cfg.Directory.Store, args[0])

	o.RootPackage = config.Cfg.RootPackage
	o.ModelPath = config.Cfg.Directory.Model
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}

	return nil
}

// Complete completes all the required options.
func (o *StoreOptions) Complete(cmd *cobra.Command, args []string) error {
	// Read template
	cmdTemplateBytes, _ := tplFS.ReadFile("tpl/store.tpl")
	cmdTemplate = string(cmdTemplateBytes)

	return nil
}

// Run executes a new sub command using the specified options.
func (o *StoreOptions) Run(args []string) error {
	return cmdutil.GenerateGoCode(o.FilePath, cmdTemplate, o)
}
