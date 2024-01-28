package generator

import (
	"embed"
	"fmt"
)

type Tmpl string

var (
	//go:embed tpl
	tplFS embed.FS
)

const (
	TmplCmd        Tmpl = "cmd"
	TmplModel      Tmpl = "model"
	TmplStore      Tmpl = "store"
	TmplRequest    Tmpl = "request"
	TmplBiz        Tmpl = "biz"
	TmplController Tmpl = "controller"
	TmplMiddleware Tmpl = "middleware"
	TmplJob        Tmpl = "job"
	TmplMigration  Tmpl = "migration"
	TmplSeeder     Tmpl = "seeder"
)

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
