package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/cmd/create"
	"github.com/bingo-project/bingoctl/pkg/cmd/db"
	"github.com/bingo-project/bingoctl/pkg/cmd/gen"
	makecmd "github.com/bingo-project/bingoctl/pkg/cmd/make"
	"github.com/bingo-project/bingoctl/pkg/cmd/migrate"
	"github.com/bingo-project/bingoctl/pkg/cmd/version"
	"github.com/bingo-project/bingoctl/pkg/config"
)

var (
	CfgFile string
)

func NewDefaultBingoCommand() *cobra.Command {
	return NewBingoCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewBingoCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var cmds = &cobra.Command{
		Use:   "bingo",
		Short: "Scaffold and code generator for Bingo framework",
		Long: `Bingo is a scaffold and code generation tool for Go,
used to quickly create and develop applications based on the Bingo framework.`,
		Run: runHelp,
	}

	// Load config
	cobra.OnInitialize(initConfig)

	// Config file
	cmds.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "The path to the configuration file. Empty string for no configuration file.")

	// Add commands
	cmds.AddCommand(version.NewCmdVersion())
	cmds.AddCommand(makecmd.NewCmdMake())
	cmds.AddCommand(create.NewCmdCreate())
	cmds.AddCommand(gen.NewCmdGen())
	cmds.AddCommand(migrate.NewCmdMigrateWithRunner())
	cmds.AddCommand(db.NewCmdDB())

	return cmds
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.LoadConfig(CfgFile, &config.Cfg)
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
