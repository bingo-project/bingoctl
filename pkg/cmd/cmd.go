package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/cmd/create"
	"github.com/bingo-project/bingoctl/pkg/cmd/gen"
	makecmd "github.com/bingo-project/bingoctl/pkg/cmd/make"
	"github.com/bingo-project/bingoctl/pkg/cmd/migrate"
	"github.com/bingo-project/bingoctl/pkg/cmd/version"
	"github.com/bingo-project/bingoctl/pkg/config"
)

var (
	CfgFile string
)

func NewDefaultBingoCtlCommand() *cobra.Command {
	return NewBingoCtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewBingoCtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var cmds = &cobra.Command{
		Use:   "bingoctl",
		Short: "bingoctl is the bingoctl startup client",
		Long:  `bingoctl is the client side tool for bingoctl startup.`,
		Run:   runHelp,
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
	cmds.AddCommand(migrate.NewCmdMigrate())

	return cmds
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.LoadConfig(CfgFile, &config.Cfg)
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
