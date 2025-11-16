package main

import (
	"github.com/spf13/cobra"

	"{{.RootPackage}}/internal/{{.ServiceName}}"
)

func main() {
	command := {{.ServiceName}}.NewAppCommand()
	cobra.CheckErr(command.Execute())
}