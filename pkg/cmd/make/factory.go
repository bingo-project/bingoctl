package make

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/mgutz/ansi"
	"gorm.io/gen"

	"github.com/bingo-project/bingoctl/pkg/config"
	cmdutil "github.com/bingo-project/bingoctl/pkg/util"
)

type Tmpl string

const (
	TmplCmd        Tmpl = "cmd"
	TmplModel      Tmpl = "model"
	TmplStore      Tmpl = "store"
	TmplRequest    Tmpl = "request"
	TmplBiz        Tmpl = "biz"
	TmplController Tmpl = "controller"
	TmplMiddleware Tmpl = "middleware"
	TmplJob        Tmpl = "job"
)

type Options struct {
	// Code template
	Name              string
	Description       string
	FilePath          string
	Directory         string
	CodeTemplate      string
	InterfaceTemplate string
	RegisterTemplate  string

	// Code attributes - variable
	PackageName        string
	StructName         string
	StructNamePlural   string
	VariableName       string
	VariableNameSnake  string
	VariableNamePlural string
	TableName          string
	ModelName          string

	// Code attributes - import path
	RootPackage string
	BizPath     string
	StorePath   string
	RequestPath string
	ModelPath   string

	// Generate by gorm.gen
	Table         string
	FieldTemplate string
	Fields        string
	MainFields    string
}

func (o *Options) GenerateCode(tmpl, path string) error {
	dir := GetMapDirectory(tmpl)

	o.SetName(tmpl)
	o.ReadCodeTemplates()
	o.GenerateAttributes(dir, path)

	// Generate from db table.
	dbTemplates := []Tmpl{TmplModel, TmplRequest}
	if slices.Contains(dbTemplates, Tmpl(o.Name)) && o.Table != "" {
		_ = o.GetFieldsFromDB()
	}

	err := cmdutil.GenerateCode(o.FilePath, o.CodeTemplate, o.Name, o)
	if err != nil {
		return err
	}

	if o.Name == string(TmplStore) {
		err = o.Register(config.Cfg.Registries.Store, o.InterfaceTemplate, o.RegisterTemplate)
		if err != nil {
			return err
		}
	}

	if o.Name == string(TmplBiz) {
		err = o.Register(config.Cfg.Registries.Biz, o.InterfaceTemplate, o.RegisterTemplate)
		if err != nil {
			return err
		}
	}

	// Format code
	cmd := exec.Command("gofmt", "-w", o.FilePath)
	_ = cmd.Run()

	return nil
}

func GetMapDirectory(tmpl string) (dir string) {
	dir = config.Cfg.Directory.CMD
	if tmpl == string(TmplModel) {
		dir = config.Cfg.Directory.Model
	}
	if tmpl == string(TmplStore) {
		dir = config.Cfg.Directory.Store
	}
	if tmpl == string(TmplRequest) {
		dir = config.Cfg.Directory.Request
	}
	if tmpl == string(TmplBiz) {
		dir = config.Cfg.Directory.Biz
	}
	if tmpl == string(TmplController) {
		dir = config.Cfg.Directory.Controller
	}
	if tmpl == string(TmplMiddleware) {
		dir = config.Cfg.Directory.Middleware
	}
	if tmpl == string(TmplJob) {
		dir = config.Cfg.Directory.Job
	}

	return
}

func (o *Options) SetName(name string) *Options {
	o.Name = name

	return o
}

func (o *Options) ReSetDirectory() *Options {
	o.Directory = ""
	o.PackageName = ""

	return o
}

func (o *Options) ReadCodeTemplates() *Options {
	// Read template
	codeTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s.tpl", o.Name))
	o.CodeTemplate = string(codeTemplateBytes)

	// Read interface template
	interfaceTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s_interface.tpl", o.Name))
	o.InterfaceTemplate = string(interfaceTemplateBytes)

	// Ream registry template
	registerTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s_registry.tpl", o.Name))
	o.RegisterTemplate = string(registerTemplateBytes)

	// Read field template
	fieldTemplateBytes, _ := tplFS.ReadFile(fmt.Sprintf("tpl/%s_field.tpl", o.Name))
	o.FieldTemplate = string(fieldTemplateBytes)

	return o
}

func (o *Options) GenerateAttributes(directory string, path string) *Options {
	// Set code attributes
	o.RootPackage = config.Cfg.RootPackage
	o.BizPath = config.Cfg.Directory.Biz
	o.StorePath = config.Cfg.Directory.Store
	o.RequestPath = config.Cfg.Directory.Request
	o.ModelPath = config.Cfg.Directory.Model

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

	// Flags: Model name
	if o.ModelName == "" {
		o.ModelName = o.StructName
	}

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

func (o *Options) Register(registry config.Registry, interfaceTemplate, registerTemplate string) error {
	if registry.Filepath == "" {
		return nil
	}

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
		registerTemplate = strings.ReplaceAll(registerTemplate, search, replace)
	}

	content, err := os.ReadFile(registry.Filepath)
	if err != nil {
		return err
	}

	// 注册 interface
	newContent, err := RegisterInterface(registry.Interface, string(content), interfaceTemplate, registerTemplate)
	if err != nil {
		return err
	}

	err = os.WriteFile(registry.Filepath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n", ansi.Color("Registered:", "green"), registry.Filepath)

	return nil
}

func (o *Options) GetFieldsFromDB() error {
	// Generate model from table.
	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: o.Directory,

		// generate model global configuration
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(config.DB)

	// Generate struct `StructName` based on table `Table`
	meta := g.GenerateModelAs(o.Table, o.StructName)
	if len(meta.Fields) == 0 {
		return nil
	}

	gormFields := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt"}
	for _, field := range meta.Fields {
		// Comment
		var comment string
		if field.ColumnComment != "" {
			comment = " // " + field.ColumnComment
		}

		// Type
		if field.Type == "gorm.DeletedAt" && o.Name == string(TmplRequest) {
			field.Type = "*time.Time"
		}

		// Replaces
		replaces := make(map[string]string)
		replaces["{{.Name}}"] = field.Name
		replaces["{{.Type}}"] = field.Type
		replaces["{{.GORMTag}}"] = field.GORMTag.Build()
		replaces["{{.JsonTag}}"] = strcase.ToLowerCamel(field.Tag["json"])
		replaces["{{.Comment}}"] = comment

		// Replace
		fieldTemplate := o.FieldTemplate
		for search, replace := range replaces {
			fieldTemplate = strings.ReplaceAll(fieldTemplate, search, replace)
		}

		// Fields
		o.Fields += fieldTemplate + "\n"

		// Skip gorm.Model fields.
		if slices.Contains(gormFields, field.Name) {
			continue
		}

		o.MainFields += fieldTemplate + "\n"
	}

	return nil
}

func RegisterInterface(name, content, interfaceTemplate, registerTemplate string) (data string, err error) {
	// Check if already registered
	if strings.Contains(content, interfaceTemplate) {
		return "", errors.New("interface already registered: " + interfaceTemplate)
	}

	// Register
	seek := fmt.Sprintf("type %s interface {", name)
	rule := seek + `([^}]*)}`
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
	newContent = newContent + "\n" + registerTemplate

	return newContent, nil
}
