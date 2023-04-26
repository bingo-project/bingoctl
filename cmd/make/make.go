package make

import (
	"path/filepath"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/util"
)

var (
	makeExample = "make cmd"
	opt         = NewOptions()
)

type Options struct {
	Name               string
	FilePath           string
	Directory          string
	PackageName        string
	StructName         string
	StructNamePlural   string
	VariableName       string
	VariableNameSnake  string
	VariableNamePlural string
	TableName          string
}

// NewOptions returns an initialized CmdOptions instance.
func NewOptions() *Options {
	return &Options{}
}

// NewCmdMake returns new initialized instance of 'new' sub command.
func NewCmdMake() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "make COMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Generate code",
		Example:               makeExample,
		Run:                   util.DefaultSubCommandRun(),
	}

	cmd.PersistentFlags().StringVarP(&opt.Directory, "directory", "d", "", "Where to create the file.")
	cmd.PersistentFlags().StringVarP(&opt.PackageName, "package", "p", "", "Name of the package.")

	// Add subcommands
	cmd.AddCommand(NewCmdCMD())
	cmd.AddCommand(NewCmdModel())
	cmd.AddCommand(NewCmdStore())
	cmd.AddCommand(NewCmdRequest())
	cmd.AddCommand(NewCmdBiz())
	cmd.AddCommand(NewCmdController())
	cmd.AddCommand(NewCmdCrud())
	cmd.AddCommand(NewCmdMiddleware())

	return cmd
}

func (o *Options) MakeOptionsFromPath(directory string, path string) *Options {
	if o.Directory == "" {
		o.Directory = directory
	}

	arr := strings.Split(filepath.Join(o.Directory, path), "/")
	name := arr[len(arr)-1]

	o.StructName = strcase.ToCamel(name)
	o.StructNamePlural = pluralize.NewClient().Plural(o.StructName)
	o.VariableName = strcase.ToLowerCamel(o.StructName)
	o.VariableNameSnake = strcase.ToSnake(o.StructName)
	o.VariableNamePlural = pluralize.NewClient().Plural(o.VariableName)
	o.TableName = strcase.ToSnake(o.StructName)

	// Directory
	directoryArr := arr[:len(arr)-1]
	o.Directory = strings.Join(directoryArr, "/")
	o.Directory = o.Directory + "/"
	if o.PackageName == "" && len(directoryArr) > 0 {
		o.PackageName = strcase.ToSnake(directoryArr[len(directoryArr)-1])
	}

	// File path
	o.FilePath = filepath.Join(o.Directory, o.VariableNameSnake+".go")

	return o
}
