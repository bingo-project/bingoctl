package {{.ServiceName}}

import (
	"github.com/bingo-project/component-base/cli"
	"github.com/spf13/cobra"
)

// NewAppCommand creates the application command.
func NewAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "{{.ServiceName}}",
		Short: "{{.ServiceName}} service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	cli.AddConfigFlag(cmd, "{{.ServiceName}}")

	return cmd
}