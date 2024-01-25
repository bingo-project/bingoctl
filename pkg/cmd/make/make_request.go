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
	requestUsageStr = "request NAME"
)

var (
	requestUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the request command",
		bizUsageStr,
	)
)

// RequestOptions is an option struct to support 'request' sub command.
type RequestOptions struct {
	*generator.Options
}

// NewRequestOptions returns an initialized RequestOptions instance.
func NewRequestOptions() *RequestOptions {
	return &RequestOptions{
		opt,
	}
}

// NewCmdRequest returns new initialized instance of 'request' sub command.
func NewCmdRequest() *cobra.Command {
	o := NewRequestOptions()

	cmd := &cobra.Command{
		Use:                   requestUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate request code",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringVarP(&o.Table, "table", "t", "", "generate by table, example:'post'.")

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *RequestOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, requestUsageErrStr)
	}

	return nil
}

// Complete completes all the required options.
func (o *RequestOptions) Complete(cmd *cobra.Command, args []string) error {
	// Init store if generating model by tables.
	var err error
	if o.Table != "" {
		config.DB, err = db.NewMySQL(config.Cfg.MysqlOptions)
	}

	return err
}

// Run executes a new sub command using the specified options.
func (o *RequestOptions) Run(args []string) error {
	return o.GenerateCode("request", args[0])
}
