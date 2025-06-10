package main

import (
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/admserver"
)

func main() {
	command := admserver.NewAppCommand()
	cobra.CheckErr(command.Execute())
}
