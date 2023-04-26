package main

import (
	"os"

	"{[.RootPackage]}/internal/watcher"
)

func main() {
	command := watcher.NewWatcherCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
