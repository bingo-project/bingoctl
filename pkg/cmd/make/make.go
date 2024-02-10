package make

import (
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/generator"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

var (
	makeExample = "make cmd"
	opt         = NewOptions()
)

// NewOptions returns an initialized CmdOptions instance.
func NewOptions() *generator.Options {
	return &generator.Options{}
}

// NewCmdMake returns new initialized instance of 'new' sub command.
func NewCmdMake() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "make COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Generate code",
		Example:               makeExample,
		Run:                   cmdutil.DefaultSubCommandRun(),
	}

	cmd.PersistentFlags().StringVarP(&opt.Directory, "directory", "d", "", "Where to create the file.")
	cmd.PersistentFlags().StringVarP(&opt.PackageName, "package", "p", "", "Name of the package.")
	cmd.PersistentFlags().StringVarP(&opt.Table, "table", "t", "", "Read fields from db table.")

	// Add subcommands
	cmd.AddCommand(NewCmdCMD())
	cmd.AddCommand(NewCmdModel())
	cmd.AddCommand(NewCmdStore())
	cmd.AddCommand(NewCmdRequest())
	cmd.AddCommand(NewCmdBiz())
	cmd.AddCommand(NewCmdController())
	cmd.AddCommand(NewCmdCrud())
	cmd.AddCommand(NewCmdMiddleware())
	cmd.AddCommand(NewCmdJob())
	cmd.AddCommand(NewCmdMigration())
	cmd.AddCommand(NewCmdSeeder())

	return cmd
}
