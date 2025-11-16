// ABOUTME: make service 子命令，用于生成服务模块
// ABOUTME: 支持通过标志配置 HTTP/gRPC 服务器和业务层目录

package make

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	serviceUsageStr = "service NAME"
)

var (
	serviceUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the service command",
		serviceUsageStr,
	)
)

// ServiceOptions is an option struct to support 'service' sub command.
type ServiceOptions struct {
	*generator.Options
}

// NewServiceOptions returns an initialized ServiceOptions instance.
func NewServiceOptions() *ServiceOptions {
	return &ServiceOptions{
		Options: opt,
	}
}

// NewCmdService returns new initialized instance of 'service' sub command.
func NewCmdService() *cobra.Command {
	o := NewServiceOptions()

	cmd := &cobra.Command{
		Use:                   serviceUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Generate service code",
		Long:                  "Generate a new service module with configurable HTTP/gRPC servers and business layers.",
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().BoolVar(&o.EnableHTTP, "http", false, "Enable HTTP server")
	cmd.Flags().BoolVar(&o.EnableGRPC, "grpc", false, "Enable gRPC server")
	cmd.Flags().BoolVar(&o.WithBiz, "with-biz", false, "Generate biz layer")
	cmd.Flags().BoolVar(&o.WithStore, "with-store", false, "Generate store layer")
	cmd.Flags().BoolVar(&o.WithController, "with-controller", false, "Generate controller layer")
	cmd.Flags().BoolVar(&o.WithMiddleware, "with-middleware", false, "Generate middleware directory")
	cmd.Flags().BoolVar(&o.WithRouter, "with-router", false, "Generate router directory")

	return cmd
}

// Validate makes sure there is no discrepancy in command options.
func (o *ServiceOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, serviceUsageErrStr)
	}

	o.ServiceName = args[0]

	// Check if cmd/ and internal/ directories exist
	if _, err := os.Stat("cmd"); os.IsNotExist(err) {
		return fmt.Errorf("cmd/ directory does not exist, please run this command in a project root")
	}
	if _, err := os.Stat("internal"); os.IsNotExist(err) {
		return fmt.Errorf("internal/ directory does not exist, please run this command in a project root")
	}

	// Check if service already exists
	cmdPath := filepath.Join("cmd", o.ServiceName)
	if _, err := os.Stat(cmdPath); !os.IsNotExist(err) {
		return fmt.Errorf("service already exists: %s", cmdPath)
	}

	internalPath := filepath.Join("internal", o.ServiceName)
	if _, err := os.Stat(internalPath); !os.IsNotExist(err) {
		return fmt.Errorf("service already exists: %s", internalPath)
	}

	return nil
}

// Complete completes all the required options.
func (o *ServiceOptions) Complete(cmd *cobra.Command, args []string) error {
	// Copy flags to generator options
	o.Options.ServiceName = o.ServiceName
	return nil
}

// Run executes a new sub command using the specified options.
func (o *ServiceOptions) Run(args []string) error {
	return o.GenerateService(args[0])
}