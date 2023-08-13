package create

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/bingo-project/component-base/cli/console"
	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

const (
	createUsageStr = "create NAME"
)

var (
	//go:embed tpl
	tplFS embed.FS
	root  = "tpl"

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nNAME is a required argument for the create command",
		createUsageStr,
	)
)

// CreateOptions is an option struct to support 'create' sub command.
type CreateOptions struct {
	GoVersion    string
	TemplatePath string
	RootPackage  string
	AppName      string
	AppNameCamel string
}

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		GoVersion: cmdutil.GetGoVersion(),
	}
}

// NewCmdCreate returns new initialized instance of 'create' sub command.
func NewCmdCreate() *cobra.Command {
	o := NewCreateOptions()

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Short:                 "Create a project",
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
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmdutil.UsageErrorf(cmd, createUsageErrStr)
	}

	// Root Package & App path
	o.RootPackage = args[0]
	arr := strings.Split(o.RootPackage, "/")
	o.AppName = strings.ToLower(arr[len(arr)-1])
	o.AppNameCamel = strcase.ToCamel(o.AppName)

	// Path is exists
	exists, err := cmdutil.PathExists(o.AppName)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	// Project exists
	console.Warn("project directory already exists!")
	prompt := promptui.Prompt{
		Label:     "Overwrite",
		IsConfirm: true,
	}

	_, err = prompt.Run()
	if err != nil {
		console.Exit("skipped.")
	}

	cmdutil.Overwrite = true

	return nil
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(cmd *cobra.Command, args []string) error {
	err := os.MkdirAll(o.AppName, 0755)
	if err != nil {
		return err
	}

	return nil
}

// Run executes a new sub command using the specified options.
func (o *CreateOptions) Run(args []string) error {
	console.Info(fmt.Sprintf("Creating project %s", o.RootPackage))
	err := fs.WalkDir(tplFS, root, func(path string, d fs.DirEntry, err error) error {
		// Cannot happen
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(tplFS, path)
		if err != nil {
			return err
		}

		// Copy files
		dest := strings.Replace(path, root+"/", o.AppName+"/", 1)
		dest = strings.Replace(dest, "."+root, "", 1)
		dest = strings.Replace(dest, "{app}", o.AppName, 1)
		dest = strings.Replace(dest, "hidden", "", 1)

		// Generate code
		err = cmdutil.GenerateCode(dest, string(b), "init", o)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	console.Info("done.")

	return nil
}
