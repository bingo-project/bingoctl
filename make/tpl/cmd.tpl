package {{.PackageName}}

import (
    "github.com/spf13/cobra"
)

const (
    {{.VariableName}}UsageStr = "{{.VariableName}}"
)

// {{.StructName}}Options is an option struct to support '{{.VariableName}}' sub command.
type {{.StructName}}Options struct {
    // Options
}

// New{{.StructName}}Options returns an initialized {{.StructName}}Options instance.
func New{{.StructName}}Options() *{{.StructName}}Options {
    return &{{.StructName}}Options{}
}

// NewCmd{{.StructName}} returns new initialized instance of '{{.VariableName}}' sub command.
func NewCmd{{.StructName}}() *cobra.Command {
    o := New{{.StructName}}Options()

    cmd := &cobra.Command{
        Use:                   {{.VariableName}}UsageStr,
        DisableFlagsInUseLine: true,
        Short:                 "{{.CommandDescription}}",
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
func (o *{{.StructName}}Options) Validate(cmd *cobra.Command, args []string) error {
    return nil
}

// Complete completes all the required options.
func (o *{{.StructName}}Options) Complete(cmd *cobra.Command, args []string) error {
    return nil
}

// Run executes a new sub command using the specified options.
func (o *{{.StructName}}Options) Run(args []string) error {
    return nil
}
