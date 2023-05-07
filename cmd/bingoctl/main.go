package main

import (
	"os"

	"github.com/bingo-project/bingoctl/internal/bingoctl/cmd"
)

func main() {
	command := cmd.NewDefaultBingoCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
