package main

import (
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/watcher"
)

func main() {
	command := watcher.NewWatcherCommand()
	cobra.CheckErr(command.Execute())
}
