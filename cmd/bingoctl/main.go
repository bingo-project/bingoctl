package main

import (
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/internal/bingoctl/cmd"
)

func main() {
	command := cmd.NewDefaultBingoCtlCommand()
	cobra.CheckErr(command.Execute())
}
