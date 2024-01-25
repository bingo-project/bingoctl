package generator

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"gorm.io/gen"

	"github.com/bingo-project/bingoctl/pkg/config"
)

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
