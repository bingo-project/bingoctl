package make

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bingo-project/component-base/cli/console"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"

	"github.com/bingo-project/bingoctl/pkg/config"
	"github.com/bingo-project/bingoctl/pkg/util"
)

var (
	//go:embed tpl
	tplFS       embed.FS
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

func (o *Options) Register(registry config.Registry, interfaceTemplate, codeTemplate string) error {
	// Package
	pkg := ""
	if o.PackageName != o.Name {
		pkg = o.PackageName + "."
	}

	// Replace
	replaces := make(map[string]string)
	replaces["{{.Package}}"] = pkg
	replaces["{{.StructName}}"] = o.StructName
	replaces["{{.StructNamePlural}}"] = o.StructNamePlural
	replaces["{{.VariableName}}"] = o.VariableName
	replaces["{{.VariableNameSnake}}"] = o.VariableNameSnake

	for search, replace := range replaces {
		interfaceTemplate = strings.ReplaceAll(interfaceTemplate, search, replace)
		codeTemplate = strings.ReplaceAll(codeTemplate, search, replace)
	}

	content, err := os.ReadFile(registry.Filepath)
	if err != nil {
		return err
	}

	// 注册 interface
	newContent, err := RegisterInterface(registry.Interface, string(content), interfaceTemplate, codeTemplate)
	if err != nil {
		return err
	}

	err = os.WriteFile(registry.Filepath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf(" - Registered %s to ", o.Name)
	console.Info(registry.Filepath)

	return nil
}

func RegisterInterface(name, content, interfaceTemplate, codeTemplate string) (data string, err error) {
	seek := fmt.Sprintf("type %s interface {", name)
	rule := seek + `[a-zA-Z0-9().*\s]*}`
	reg := regexp.MustCompile(rule)

	// 根据规则提取关键信息
	results := reg.FindAllString(content, -1)
	if len(results) == 0 {
		err = errors.New("not matched")

		return
	}

	old := results[0]
	str := strings.TrimRight(old, "}")
	str = str + "\t" + interfaceTemplate + "\n}"

	newContent := strings.Replace(content, old, str, 1)
	newContent = newContent + "\n" + codeTemplate

	return newContent, nil
}
