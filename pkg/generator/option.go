package generator

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
	RootPackage  string
	BizPath      string
	StorePath    string
	RequestPath  string
	ModelPath    string
	RelativePath string

	// Generate by gorm.gen
	Table         string
	FieldTemplate string
	Fields        string
	MainFields    string

	// Migration
	TimeStr string
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
