package main

import (
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/scheduler"
)

func main() {
	command := scheduler.NewSchedulerCommand()
	cobra.CheckErr(command.Execute())
}
