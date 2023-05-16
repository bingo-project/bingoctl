package main

import (
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/{[.AppName]}ctl/cmd"
)

func main() {
	command := cmd.NewDefault{[.AppNameCamel]}CtlCommand()
	cobra.CheckErr(command.Execute())
}
