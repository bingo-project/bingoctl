package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/cmd/create"
	"github.com/bingo-project/bingoctl/cmd/gen"
	makecmd "github.com/bingo-project/bingoctl/cmd/make"
	"github.com/bingo-project/bingoctl/config"
)

func main() {
	command := NewBingoCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewBingoCtlCommand() *cobra.Command {
	var cmds = &cobra.Command{
		Use:   "bingoctl",
		Short: "bingoctl is the bingo client",
		Long:  `bingoctl is the client side tool for Bingo project.`,
		Run:   runHelp,
	}

	// Load config
	cobra.OnInitialize(initConfig)

	// Config file
	cmds.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "The path to the configuration file. Empty string for no configuration file.")

	// Add commands
	cmds.AddCommand(makecmd.NewCmdMake())
	cmds.AddCommand(create.NewCmdCreate())
	cmds.AddCommand(gen.NewCmdGen())

	return cmds
}

var (
	CfgFile string
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.LoadConfig(CfgFile, &config.Cfg)
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
