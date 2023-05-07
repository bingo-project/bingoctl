package main

import (
	"os"

	"{[.RootPackage]}/internal/{[.AppName]}ctl/cmd"
)

func main() {
	command := cmd.NewDefault{[.AppNameCamel]}CtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
