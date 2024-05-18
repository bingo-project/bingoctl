package main

import (
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/bot"
)

func main() {
	command := bot.NewBotCommand()
	cobra.CheckErr(command.Execute())
}
