package cmd

import (
	"io"
	"os"

	"github.com/bingo-project/component-base/cli/genericclioptions"
	"github.com/bingo-project/component-base/cli/templates"
	"github.com/spf13/cobra"

	"{[.RootPackage]}/internal/apiserver/bootstrap"
	"{[.RootPackage]}/internal/{[.AppName]}ctl/cmd/user"
	"{[.RootPackage]}/internal/{[.AppName]}ctl/cmd/version"
)

func NewDefault{[.AppNameCamel]}CtlCommand() *cobra.Command {
	return New{[.AppNameCamel]}CtlCommand(os.Stdin, os.Stdout, os.Stderr)
}

func New{[.AppNameCamel]}CtlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var cmds = &cobra.Command{
		Use:   "{[.AppName]}ctl",
		Short: "{[.AppName]}ctl is the {[.AppName]} startup client",
		Long:  `{[.AppName]}ctl is the client side tool for {[.AppName]} startup.`,
		Run:   runHelp,
	}

	// Load config
	cobra.OnInitialize(initConfig)

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	groups := templates.CommandGroups{
		{
			Message:  "Basic Commands:",
			Commands: []*cobra.Command{},
		},
		{
			Message: "Advanced Commands:",
			Commands: []*cobra.Command{
				user.NewCmdUser(ioStreams),
			},
		},
	}
	groups.Add(cmds)

	filters := []string{""}
	templates.ActsAsRootCommand(cmds, filters, groups...)

	// Config file
	cmds.PersistentFlags().StringVarP(&bootstrap.CfgFile, "config", "c", "", "The path to the configuration file. Empty string for no configuration file.")

	// Add commands
	cmds.AddCommand(version.NewCmdVersion(ioStreams))

	return cmds
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	bootstrap.InitConfig()
	bootstrap.Boot()
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
